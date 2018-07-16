package helpers

import (
	"encoding/json"
	"log"

	"github.com/byuoitav/common/status"
)

func GetVolume(address string) (status.Volume, error) {
	log.Printf("Getting volume for %v", address)
	parentResponse, err := getAudioInformation(address)
	if err != nil {
		return status.Volume{}, err
	}
	log.Printf("%v", parentResponse)

	var output status.Volume
	for _, outerResult := range parentResponse.Result {

		for _, result := range outerResult {

			if result.Target == "speaker" {

				output.Volume = result.Volume
			}
		}
	}
	log.Printf("Done")

	return output, nil
}

func getAudioInformation(address string) (SonyAudioResponse, error) {
	payload := SonyTVRequest{
		Params:  []map[string]interface{}{},
		Method:  "getVolumeInformation",
		Version: "1.0",
		ID:      1,
	}

	log.Printf("%+v", payload)

	resp, err := PostHTTP(address, payload, "audio")

	parentResponse := SonyAudioResponse{}

	log.Printf("%s", resp)

	err = json.Unmarshal(resp, &parentResponse)
	return parentResponse, err

}

func GetMute(address string) (status.Mute, error) {
	log.Printf("Getting mute status for %v", address)
	parentResponse, err := getAudioInformation(address)
	if err != nil {
		return status.Mute{}, err
	}
	var output status.Mute
	for _, outerResult := range parentResponse.Result {
		for _, result := range outerResult {
			if result.Target == "speaker" {
				log.Printf("local mute: %v", result.Mute)
				output.Muted = result.Mute
			}
		}
	}

	log.Printf("Done")

	return output, nil
}

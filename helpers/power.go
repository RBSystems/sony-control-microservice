package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/nerr"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
)

// SetPower sets the power status through the testing Sony API
func SetPower(address string, status bool) *nerr.E {

	// Make params which is what will contain the power on/off command
	params := make(map[string]interface{})

	//Status will be a bool since power on/off is in terms of true/false
	params["status"] = status

	// currentStatus will get the current power state of the projector
	currentStatus, err := GetPower(address)
	if err != nil {
		return nerr.Translate(err)
	}

	// Build the payload that will be sent to change the power.
	err = BuildAndSendPayload(address, "system", "setPowerStatus", params)
	if err != nil {
		return nerr.Translate(err)
	}

	//we need to wait for a little bit to let the tv finish so it doesn't override
	if currentStatus.Power == "on" && !status {
		log.L.Infof("Waiting to go to standby....")
		time.Sleep(1750 * time.Millisecond)
	} else if currentStatus.Power == "standby" && status {
		log.L.Infof("Waiting to turn on...")
		time.Sleep(1750 * time.Millisecond)
	}

	postStatus, err := GetPower(address)
	if err != nil {
		return nerr.Translate(err)
	}

	log.L.Infof("%v", postStatus)

	// There is some issue with the projector right here...
	if status && postStatus.Power != "on" {
		// do we want to retry the command
		return nerr.Create(fmt.Sprintf("Power wasn't set successfully from %v", postStatus), "Failed") //nerr.Create()
	} else if !status && postStatus.Power != "standby" {
		return nerr.Create(fmt.Sprintf("Power wasn't set to standby correctly %v", postStatus), "Failed")
	}

	return nil
}

// GetPower retrieves the power status through the testing Sony API
func GetPower(address string) (status.Power, *nerr.E) {

	// powerOutput is the status.Power JSON thingy...
	var powerOutput status.Power

	request := SonyRequest{
		Method:  "getPowerStatus",
		Version: "1.0",
		ID:      1,
		Params:  []map[string]interface{}{},
	}

	response, err := PostHTTP(address, request, "system")
	if err != nil {
		return status.Power{}, err
	}

	powerStatus := string(response)

	log.L.Debugf("Device returned: %s", powerStatus)

	if strings.Contains(powerStatus, "active") {
		powerOutput.Power = "on"
	} else if strings.Contains(powerStatus, "standby") {
		powerOutput.Power = "standby"
	} else {
		return status.Power{}, nerr.Translate(err).Addf("There was an error getting power status")
	}

	return powerOutput, nil
}

package main

import (
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/sony-control-microservice/handlers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8007"
	router := echo.New()

	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	//functionality endpoints
	secure.GET("/:address/power/on", handlers.PowerOn)
	secure.GET("/:address/power/standby", handlers.Standby)
	secure.GET("/:address/input/:port", handlers.SwitchInput)
	secure.GET("/:address/volume/set/:value", handlers.SetVolume)
	secure.GET("/:address/volume/mute", handlers.VolumeMute)
	secure.GET("/:address/volume/unmute", handlers.VolumeUnmute)
	secure.GET("/:address/display/blank", handlers.BlankDisplay)
	secure.GET("/:address/display/unblank", handlers.UnblankDisplay)

	// Web API functionality endpoints

	//status endpoints
	secure.GET("/:address/power/status", handlers.GetPower)
	secure.GET("/:address/input/current", handlers.GetInput)
	secure.GET("/:address/input/list", handlers.GetInputList)
	secure.GET("/:address/volume/level", handlers.GetVolume)
	secure.GET("/:address/volume/mute/status", handlers.GetMute)
	secure.GET("/:address/display/status", handlers.GetBlank)

	// Web API status endpoints
	secure.GET("/:address/getPowerStatus/status", handlers.GetPowerAPI)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}

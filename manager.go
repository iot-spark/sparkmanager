package main

import (
	"github.com/iot-spark/sparkmanager/controllers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

func Manager() {
	// init router
	router := httprouter.New()

	// Init device controller
	DevContr := controllers.NewDeviceController(*DbFile)

	// Get available devices
	router.GET("/devices", DevContr.GetDevices)

	// Register a new device
	router.POST("/devices/add", DevContr.AddDevice)

	// Remove a particular device
	router.DELETE("/devices/delete/:id", DevContr.RemoveDevice)

	log.Fatal(
		http.ListenAndServe(
			strings.Join([]string{"", *ManagerPort}, ":"),
			router,
		),
	)

}

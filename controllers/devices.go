package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/iot-spark/sparkmanager/models"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type (
	DeviceController struct {
		dm *models.DeviceModel
	}

	Response struct {
		Status bool            `json:"status"`
		Data   []models.Device `json:"data,omitempty"`
		Error  string          `json:"message,omitempty"`
	}
)

func NewDeviceController(dbFile string) *DeviceController {
	// Init DeviceModel
	dm := models.NewDeviceModel(dbFile)

	return &DeviceController{dm}
}

/**
 * Retrieves all devices
 * @return json
 */
func (c *DeviceController) GetDevices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := Response{}

	devices, status := c.dm.GetAllDevices()
	response.Status = status
	response.Data = devices

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	res := InitResponse(&response)
	fmt.Fprintf(w, "%s", res)
}

/**
 * Adds a new device
 */
func (c *DeviceController) AddDevice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := Response{}

	device, status := c.dm.AddDevice(r.Body)
	if status {
		// TODO:
		// generate psk_file
		// send SIGHUP siganl to mqtt to reload config gracefully
		response.Status = status
		response.Data = device
	} else {
		response.Error = string("Device registration error!")
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	res := InitResponse(&response)
	fmt.Fprintf(w, "%s", res)
}

/**
 * Removes an existing device
 */
func (c *DeviceController) RemoveDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	status := 400
	id := p.ByName("id")
	res := c.dm.DeleteDevice(id)
	if res {
		status = 200
	}
	w.WriteHeader(status)
}

func InitResponse(res *Response) []byte {

	response, err := json.Marshal(&res)
	if err != nil {
		log.Printf("ENCODE RESPONSE: Json encode error: %#v\n", err)
		return []byte("")
	}

	return response
}

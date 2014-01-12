package rest

import (
	"github.com/emicklei/go-restful"
	devicePkg "github.com/gopns/gopns/device"
	"log"
	"net/http"
)

type DeviceService struct {
	DeviceManager devicePkg.DeviceManager
}

func (serv *DeviceService) Register(container *restful.Container, rootPath string) {
	ws := new(restful.WebService)
	ws.
		Path(rootPath + "/device").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(serv.getDevices).
		// docs
		Doc("list devices").
		Param(ws.QueryParameter("cursor", "the cursor for fetching the next set of devices").DataType("string")).
		Writes(devicePkg.DeviceList{}))

	ws.Route(ws.POST("/").To(serv.registerDevice).
		// docs
		Doc("register a new device").
		Reads(devicePkg.DeviceRegistration{}))

	ws.Route(ws.GET("/{deviceAlias}").To(serv.getDevice).
		// docs
		Doc("get device by alias").
		Param(ws.PathParameter("deviceAlias", "the registered device alias").DataType("string")).
		Writes(devicePkg.Device{}))

	container.Add(ws)
}

func (serv *DeviceService) getDevice(request *restful.Request, response *restful.Response) {
	alias := request.PathParameter("deviceAlias")
	err, device_ := serv.DeviceManager.GetDevice(alias)
	if err != nil || device_ == nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Device not found.")
		return
	}
	response.WriteEntity(*device_)

}

func (serv *DeviceService) getDevices(request *restful.Request, response *restful.Response) {

	cursor := request.QueryParameter("cursor")

	err, deviceList := serv.DeviceManager.GetDevices(cursor)
	if err != nil || deviceList == nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "No device registered.")
		return
	}
	log.Printf("Devices found: %v\n", *deviceList)
	response.WriteEntity(*deviceList)
}

func (serv *DeviceService) registerDevice(request *restful.Request, response *restful.Response) {

	deviceR := new(devicePkg.DeviceRegistration)
	err := request.ReadEntity(deviceR)
	if err != nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Invalid device")
		return
	}

	// ToDo validate device

	err, _ = serv.DeviceManager.RegisterDevice(*deviceR)
	if err != nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

}

func (serv *DeviceService) addTags(request *restful.Request, response *restful.Response) {

	return
}

func (serv *DeviceService) deleteDevice(request *restful.Request, response *restful.Response) {

	return
}

func (serv *DeviceService) deleteTag(request *restful.Request, response *restful.Response) {

	return
}

func (serv *DeviceService) deleteArn(request *restful.Request, response *restful.Response) {

	return
}

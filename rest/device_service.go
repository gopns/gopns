package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/device"
	"github.com/gopns/gopns/exception"
	"github.com/gopns/gopns/model"
	"log"
)

type DeviceService struct {
	DeviceManager device.DeviceManager
}

func (serv *DeviceService) Register(container *restful.Container, rootPath string) {
	ws := new(restful.WebService)
	ws.
		Path(rootPath + "/apps/{appId}/devices").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").
		Filter(NewTimingFilter("get-devices")).
		To(serv.getDevices).
		// docs
		Doc("list all devices for an app").
		Param(ws.PathParameter("appId", "the application id").DataType("string")).
		Param(ws.QueryParameter("cursor", "the cursor for fetching the next set of devices").DataType("string")).
		Writes(modelview.PaginatedListView{}))

	ws.Route(ws.POST("/").
		Filter(NewTimingFilter("register-device")).
		To(serv.registerDevice).
		// docs
		Doc("register a new device").
		Param(ws.PathParameter("appId", "the application id").DataType("string")).
		Reads(modelview.RegisterDeviceView{}))

	ws.Route(ws.GET("/{deviceId}").
		Filter(NewTimingFilter("get-device")).
		To(serv.getDevice).
		// docs
		Doc("get device by id").
		Param(ws.PathParameter("appId", "the application id").DataType("string")).
		Param(ws.PathParameter("deviceId", "the device id").DataType("string")).
		Writes(model.Device{}))

	ws.Route(ws.PUT("/{deviceId}").
		Filter(NewTimingFilter("put-device")).
		To(serv.addTags).
		// docs
		Doc("Add tag to device").
		Param(ws.PathParameter("appId", "the application id").DataType("string")).
		Param(ws.PathParameter("deviceId", "the registered device alias").DataType("string")))

	container.Add(ws)
}

func (serv *DeviceService) getDevice(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("deviceId")
	if id == nil || id == "" {
		panic(exception.BadRequestException("No device Id specified."))
	}
	device, err := serv.DeviceManager.GetDevice(id)
	exception.ConditionalThrowInternalServerErrorException(err)
	if device == nil {
		log.Printf("Device not found for id %s", alias)
		panic(exception.NotFoundException("Device Not Found"))
	}

	//convert device to device view
	deviceView := modelview.ConvertToDeviceView(device)
	response.WriteEntity(*deviceView)

}

func (serv *DeviceService) getDevices(request *restful.Request, response *restful.Response) {

	cursor := request.QueryParameter("cursor")

	err, deviceList := serv.DeviceManager.GetDevices(cursor)
	exception.ConditionalThrowInternalServerErrorException(err)

	log.Printf("Devices found: %v\n", *deviceList)
	response.WriteEntity(*deviceList)
}

func (serv *DeviceService) registerDevice(request *restful.Request, response *restful.Response) {

	deviceR := new(model.DeviceRegistration)
	err := request.ReadEntity(deviceR)

	exception.ConditionalThrowBadRequestException(err)

	// ToDo validate device

	err, _ = serv.DeviceManager.RegisterDevice(*deviceR)
	exception.ConditionalThrowInternalServerErrorException(err)

}

func (serv *DeviceService) addTags(request *restful.Request, response *restful.Response) {
	panic(exception.NotImplemented("Not Implemented"))
	return
}

func (serv *DeviceService) deleteDevice(request *restful.Request, response *restful.Response) {
	panic(exception.NotImplemented("Not Implemented"))
	return
}

func (serv *DeviceService) deleteTag(request *restful.Request, response *restful.Response) {
	panic(exception.NotImplemented("Not Implemented"))
	return
}

func (serv *DeviceService) deleteArn(request *restful.Request, response *restful.Response) {
	panic(exception.NotImplemented("Not Implemented"))
	return
}

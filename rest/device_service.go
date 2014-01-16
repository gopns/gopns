package rest

import (
	"github.com/emicklei/go-restful"
	devicePkg "github.com/gopns/gopns/device"
	"github.com/gopns/gopns/exception"
	"log"
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

	ws.Route(ws.GET("/").
		Filter(NewTimingFilter("get-devices")).
		To(serv.getDevices).
		// docs
		Doc("list devices").
		Param(ws.QueryParameter("cursor", "the cursor for fetching the next set of devices").DataType("string")).
		Writes(devicePkg.DeviceList{}))

	ws.Route(ws.POST("/").
		Filter(NewTimingFilter("post-device")).
		To(serv.registerDevice).
		// docs
		Doc("register a new device").
		Reads(devicePkg.DeviceRegistration{}))

	ws.Route(ws.GET("/{deviceAlias}").
		Filter(NewTimingFilter("get-device")).
		To(serv.getDevice).
		// docs
		Doc("get device by alias").
		Param(ws.PathParameter("deviceAlias", "the registered device alias").DataType("string")).
		Writes(devicePkg.Device{}))

	ws.Route(ws.POST("/{deviceAlias}/tag/{tag}").
		Filter(NewTimingFilter("add-tag")).
		To(serv.addTags).
		// docs
		Doc("Add tag to device").
		Param(ws.PathParameter("deviceAlias", "the registered device alias").DataType("string")).
		Param(ws.PathParameter("tag", "The tag to add").DataType("string")))

	container.Add(ws)
}

func (serv *DeviceService) getDevice(request *restful.Request, response *restful.Response) {
	alias := request.PathParameter("deviceAlias")
	err, device_ := serv.DeviceManager.GetDevice(alias)
	exception.ConditionalThrowNotFoundException(err)
	if device_ == nil {
		log.Printf("Device not found for id %s", alias)
		panic(exception.NotFoundException("Device Not Found"))
	}
	response.WriteEntity(*device_)

}

func (serv *DeviceService) getDevices(request *restful.Request, response *restful.Response) {

	cursor := request.QueryParameter("cursor")

	err, deviceList := serv.DeviceManager.GetDevices(cursor)
	exception.ConditionalThrowInternalServerErrorException(err)

	log.Printf("Devices found: %v\n", *deviceList)
	response.WriteEntity(*deviceList)
}

func (serv *DeviceService) registerDevice(request *restful.Request, response *restful.Response) {

	deviceR := new(devicePkg.DeviceRegistration)
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

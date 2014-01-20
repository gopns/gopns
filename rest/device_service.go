package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/access"
	"github.com/gopns/gopns/exception"
	"github.com/gopns/gopns/model"
	"github.com/gopns/gopns/modelview"
	"log"
)

type DeviceService struct {
	DeviceManager access.DeviceManager
}

func (serv *DeviceService) Register(container *restful.Container, rootPath string) {
	ws := new(restful.WebService)
	ws.
		Path(rootPath + "/apps/{appId}/devices").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").
		Filter(NewTimingFilter("list-devices")).
		To(serv.listDevices).
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
		Reads(modelview.DeviceRegisterView{}))

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
		To(serv.updateDevice).
		// docs
		Doc("Add tag to device").
		Param(ws.PathParameter("appId", "the application id").DataType("string")).
		Param(ws.PathParameter("deviceId", "the registered device alias").DataType("string")))

	container.Add(ws)
}

func (serv *DeviceService) getDevice(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("deviceId")
	if id == "" {
		panic(exception.BadRequestException("No device Id specified."))
	}
	device, err := serv.DeviceManager.GetDevice(id)
	exception.ConditionalThrowInternalServerErrorException(err)
	if device == nil {
		log.Printf("Device not found for id %s", id)
		panic(exception.NotFoundException("Device Not Found"))
	}

	//convert device to device view
	deviceView := modelview.FromDevice(*device)
	response.WriteEntity(*deviceView)

}

func (serv *DeviceService) updateDevice(request *restful.Request, response *restful.Response) {

	device := new(model.Device)
	err := request.ReadEntity(device)

	exception.ConditionalThrowBadRequestException(err)

	// ToDo validate device

	err = serv.DeviceManager.PutDevice(*device)
	exception.ConditionalThrowInternalServerErrorException(err)

}

func (serv *DeviceService) listDevices(request *restful.Request, response *restful.Response) {

	appId := request.PathParameter("appId")
	cursor := request.QueryParameter("cursor")

	devices, newCursor, err := serv.DeviceManager.ListAppDevices(appId, cursor)
	exception.ConditionalThrowInternalServerErrorException(err)
	log.Printf("Devices found: %v\n", devices)

	paginatedList := modelview.NewPaginatedDeviceListView(*devices, newCursor)
	response.WriteEntity(*paginatedList)
}

func (serv *DeviceService) registerDevice(request *restful.Request, response *restful.Response) {

	deviceR := new(modelview.DeviceRegisterView)
	err := request.ReadEntity(deviceR)

	exception.ConditionalThrowBadRequestException(err)

	dv, err := deviceR.ToDevice()
	exception.ConditionalThrowBadRequestException(err)

	// ToDo validate device

	_, err = serv.DeviceManager.RegisterDevice(*dv)
	exception.ConditionalThrowInternalServerErrorException(err)

}

func (serv *DeviceService) deleteDevice(request *restful.Request, response *restful.Response) {
	panic(exception.NotImplemented("Not Implemented"))
	return
}

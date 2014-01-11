package rest

import (
	"code.google.com/p/gorest"
	devicePkg "github.com/gopns/gopns/device"
	"github.com/gopns/gopns/rest/restutil"
)

type DeviceService struct {
	DeviceManager devicePkg.DeviceManager

	//Service level config
	gorest.RestService `root:"/rest/device/" consumes:"application/json" produces:"application/json"`

	registerDevice gorest.EndPoint `method:"POST" path:"/" postdata:"DeviceRegistration"`
	addTags        gorest.EndPoint `method:"POST" path:"/{deviceAlias:string}/tags/" postdata:"[]string"`
	getDevice      gorest.EndPoint `method:"GET" path:"/{deviceAlias:string}" output:"Device"`
	getDevices     gorest.EndPoint `method:"GET" path:"/?{cursor:string}" output:"DeviceList"`
	deleteTag      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/tag/{tag:string}"`
	deleteArn      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/arn/{arn:string}"`
	deleteDevice   gorest.EndPoint `method:"DELETE" path:"/{alias:string}"`
}

func (serv DeviceService) GetDevice(deviceAlias string) devicePkg.Device {
	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)
	err, device := serv.DeviceManager.GetDevice(deviceAlias)
	restutil.CheckError(err, restError, 500)
	return *device

}

func (serv DeviceService) GetDevices(cursor string) devicePkg.DeviceList {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	err, deviceList := serv.DeviceManager.GetDevices(cursor)
	restutil.CheckError(err, restError, 500)
	return *deviceList
}

func (serv DeviceService) RegisterDevice(device devicePkg.DeviceRegistration) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)
	err, code := serv.DeviceManager.RegisterDevice(device)
	restutil.CheckError(err, restError, code)

	return
}

func (serv DeviceService) AddTags(tags []string, deviceAlias string) {

	return
}

func (serv DeviceService) DeleteDevice(deviceAlias string) {

	return
}

func (serv DeviceService) DeleteTag(deviceAlias string, tag string) {

	return
}

func (serv DeviceService) DeleteArn(deviceAlias string, arn string) {

	return
}

package device

import (
	"code.google.com/p/gorest"
)

type DeviceService struct {
	//Service level config
	gorest.RestService `root:"/rest/device/" consumes:"application/json" produces:"application/json"`

	registerDevice gorest.EndPoint `method:"POST" path:"/" postdata:"DeviceRegistration"`
	addTags        gorest.EndPoint `method:"POST" path:"/{deviceAlias:string}/tags/" postdata:"[]string"`
	getDevice      gorest.EndPoint `method:"GET" path:"/{deviceAlias:string}" output:"Device"`
	getDevices     gorest.EndPoint `method:"GET" path:"/" output:"[]Device"`
	deleteTag      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/tag/{tag:string}"`
	deleteArn      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/arn/{arn:string}"`
}

func (serv DeviceService) GetDevice(deviceAlias string) Device {

	return Device{deviceAlias, "EN_US", []string{"Arn1", "Arn2"}, []string{"Whale"}}
}

func (serv DeviceService) GetDevices() []Device {

	return []Device{Device{"Alias1", "EN_US", []string{"Arn1", "Arn2"}, []string{"Whale"}}, Device{"Alias2", "EN_CA", []string{"Arn3", "Arn4"}, []string{"Minnow"}}}
}

func (serv DeviceService) RegisterDevice(device DeviceRegistration) {

	return
}

func (serv DeviceService) AddTags(tags []string, deviceAlias string) {

	return
}

func (serv DeviceService) DeleteTag(deviceAlias string, tag string) {

	return
}

func (serv DeviceService) DeleteArn(deviceAlias string, arn string) {

	return
}

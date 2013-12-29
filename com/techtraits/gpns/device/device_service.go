package device

import (
	"code.google.com/p/gorest"
)

type DeviceService struct {
	//Service level config
	gorest.RestService `root:"/device/" consumes:"application/json" produces:"application/json"`

	registerDevice gorest.EndPoint `method:"POST" path:"/" postdata:"Device"`
	getDevice      gorest.EndPoint `method:"GET" path:"/{deviceAlias:string}" output:"Device"`
	getDevices     gorest.EndPoint `method:"GET" path:"/" output:"[]Device"`
}

func (serv DeviceService) GetDevice(deviceAlias string) Device {

	return Device{deviceAlias, "DeviceID", "Arn", "IOS", "EN_US", []string{"Whale"}}
}

func (serv DeviceService) GetDevices() []Device {

	return []Device{Device{"Alias1", "DeviceID1", "Arn1", "IOS", "EN_US", []string{"Whale"}}, Device{"Alias2", "DeviceID2", "Arn2", "ANDROID", "EN_CA", []string{"Minnow"}}}
}

func (serv DeviceService) RegisterDevice(device Device) {

	return
}

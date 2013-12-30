package device

import (
	"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws/sns"
	"github.com/usmanismail/gpns/com/techtraits/gpns/rest/restutil"
)

type DeviceService struct {

	//Service level config
	gorest.RestService `root:"/rest/device/" consumes:"application/json" produces:"application/json"`

	registerDevice gorest.EndPoint `method:"POST" path:"/" postdata:"DeviceRegistration"`
	addTags        gorest.EndPoint `method:"POST" path:"/{deviceAlias:string}/tags/" postdata:"[]string"`
	getDevice      gorest.EndPoint `method:"GET" path:"/{deviceAlias:string}" output:"Device"`
	getDevices     gorest.EndPoint `method:"GET" path:"/?{cursor:string}" output:"DeviceList"`
	deleteTag      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/tag/{tag:string}"`
	deleteArn      gorest.EndPoint `method:"DELETE" path:"/{alias:string}/arn/{arn:string}"`
}

func (serv DeviceService) GetDevice(deviceAlias string) Device {

	return Device{deviceAlias, "EN_US", []string{"Arn1", "Arn2"}, []string{"Whale"}}
}

func (serv DeviceService) GetDevices(cursor string) DeviceList {

	return DeviceList{[]Device{Device{"Alias1", "EN_US", []string{"Arn1", "Arn2"}, []string{"Whale"}}, Device{"Alias2", "EN_CA", []string{"Arn3", "Arn4"}, []string{"Minnow"}}}, "cursor"}
}

func (serv DeviceService) RegisterDevice(device DeviceRegistration) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	err := device.ValidateLocale()
	restutil.CheckError(err, restError, 400)
	//TODO Register with Database

	sns.RegistrarInstance().RegisterDevice(device.PlatformApp, device.Id, formatTags(device.Locale, device.Alias, device.Tags))

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

func formatTags(locale string, alias string, tags []string) string {

	tagString := alias
	tagString = tagString + "," + locale
	for _, tag := range tags {
		tagString = tagString + "," + tag
	}
	return tagString
}

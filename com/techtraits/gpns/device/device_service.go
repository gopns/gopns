package device

import (
	"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws/dynamodb"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws/sns"
	config "github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
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
	deleteDevice   gorest.EndPoint `method:"DELETE" path:"/{alias:string}"`
}

func (serv DeviceService) GetDevice(deviceAlias string) Device {
	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: deviceAlias}
	getItemRequest := dynamodb.GetItemRequest{Key: key, TableName: config.AWSConfigInstance().DynamoTable()}

	item, err := dynamodb.GetItem(
		getItemRequest,
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())
	restutil.CheckError(err, restError, 500)

	return Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["tags"].SS}
}

func (serv DeviceService) GetDevices(cursor string) DeviceList {

	return DeviceList{[]Device{Device{"Alias1", "EN_US", []string{"Arn1", "Arn2"}, []string{"Whale"}}, Device{"Alias2", "EN_CA", []string{"Arn3", "Arn4"}, []string{"Minnow"}}}, "cursor"}
}

func (serv DeviceService) RegisterDevice(device DeviceRegistration) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	err := device.ValidateLocale()
	restutil.CheckError(err, restError, 400)

	arn, err := sns.RegisterDevice(
		device.Id,
		formatTags(device.Locale, device.Alias, device.Tags),
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region(),
		config.AWSConfigInstance().PlatformApps()[device.PlatformApp].Arn())
	restutil.CheckError(err, restError, 400)

	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: device.Alias}
	attributeUpdates := make(map[string]dynamodb.AttributeUpdate)
	attributeUpdates["arns"] = dynamodb.AttributeUpdate{"ADD", dynamodb.Attribute{SS: []string{arn}}}
	attributeUpdates["locale"] = dynamodb.AttributeUpdate{"PUT", dynamodb.Attribute{S: device.Locale}}
	attributeUpdates["tags"] = dynamodb.AttributeUpdate{"ADD", dynamodb.Attribute{SS: device.Tags}}

	updateItemRequest := dynamodb.UpdateItemRequest{
		Key:              key,
		AttributeUpdates: attributeUpdates,
		TableName:        config.AWSConfigInstance().DynamoTable()}
	err = dynamodb.UpdateItem(
		updateItemRequest,
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())
	restutil.CheckError(err, restError, 500)

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

func formatTags(locale string, alias string, tags []string) string {

	tagString := alias
	tagString = tagString + "," + locale
	for _, tag := range tags {
		tagString = tagString + "," + tag
	}
	return tagString
}

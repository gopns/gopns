package rest

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/aws/dynamodb"
	devicePkg "github.com/gopns/gopns/device"
	config "github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/rest/restutil"
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

func (serv DeviceService) GetDevice(deviceAlias string) devicePkg.Device {
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

	return devicePkg.Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
}

func (serv DeviceService) GetDevices(cursor string) devicePkg.DeviceList {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	var startKey map[string]dynamodb.Attribute
	if len(cursor) > 0 {
		startKey = make(map[string]dynamodb.Attribute)
		startKey["alias"] = dynamodb.Attribute{S: cursor}
	}

	scanRequest := dynamodb.ScanRequest{ExclusiveStartKey: startKey, TableName: config.AWSConfigInstance().DynamoTable(), Limit: 1000}

	response, err := dynamodb.ScanForItems(
		scanRequest,
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())
	restutil.CheckError(err, restError, 500)

	cursor = ""
	if len(response.LastEvaluatedKey) != 0 {
		cursor = response.LastEvaluatedKey["alias"].S
	}

	return devicePkg.DeviceList{convertToDevices(response.Items), cursor}
}

func (serv DeviceService) RegisterDevice(device devicePkg.DeviceRegistration) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)
	err, code := devicePkg.RegistrarInstance().RegisterDevice(device)
	restutil.CheckError(err, restError, code)

	return
}

func convertToDevices(items []map[string]dynamodb.Attribute) []devicePkg.Device {

	devices := make([]devicePkg.Device, 0, 0)
	for _, item := range items {
		device := devicePkg.Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
		devices = append(devices, device)
	}

	return devices
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

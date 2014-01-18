package device

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/metrics"
	"github.com/gopns/gopns/model"
	"github.com/gopns/gopns/modelview"
)

type DeviceManager interface {
	RegisterDevice(device model.Device) (int, error)
	GetDevice(deviceId string) (*model.Device, error)
	ListAppDevices(cursor string) (devices *[]model.Device, newCursor string, err error)
	PutDevice(device model.Device) error
}

type DefaultDeviceManager struct {
	SnsClient    sns.SNSClient
	DynamoClient dynamodb.DynamoClient
	PlatformApps map[string]map[string]string
	DeviceTable  string
}

func New(snsClient sns.SNSClient, dynamoClient dynamodb.DynamoClient, deviceTable string,
	platformApps map[string]map[string]string) DeviceManager {
	deviceManagerInstance := &DefaultDeviceManager{
		snsClient,
		dynamoClient,
		platformApps,
		deviceTable}
	return deviceManagerInstance
}

func (this *DefaultDeviceManager) RegisterDevice(device modelview.DeviceRegistration) (error, int) {
	callMeter, errorMeter := metrics.GetCallMeters("device_manager.register_device")
	callMeter.Mark(1)

	err := model.ValidateLocale(device.Local())
	if err != nil {
		errorMeter.Mark(1)
		return err, 400
	}

	arn, err := this.SnsClient.RegisterDevice(
		device.Id,
		formatTags(device.Locale, device.Alias, device.Tags),
		this.PlatformApps[device.PlatformApp]["Arn"])

	if err != nil {
		errorMeter.Mark(1)
		return err, 400
	}

	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: device.Alias}
	attributeUpdates := make(map[string]dynamodb.AttributeUpdate)
	attributeUpdates["arns"] = dynamodb.AttributeUpdate{"ADD", dynamodb.Attribute{SS: []string{arn}}}
	attributeUpdates["locale"] = dynamodb.AttributeUpdate{"PUT", dynamodb.Attribute{S: device.Locale}}
	attributeUpdates["tags"] = dynamodb.AttributeUpdate{"ADD", dynamodb.Attribute{SS: device.Tags}}
	attributeUpdates["platform"] = dynamodb.AttributeUpdate{"PUT", dynamodb.Attribute{S: device.PlatformApp}}

	updateItemRequest := dynamodb.UpdateItemRequest{
		Key:              key,
		AttributeUpdates: attributeUpdates,
		TableName:        this.DeviceTable}
	err = this.DynamoClient.UpdateItem(updateItemRequest)

	if err != nil {
		errorMeter.Mark(1)
		return err, 500
	}

	return nil, 0
}

func (this *DefaultDeviceManager) GetDevice(deviceAlias string) (error, *model.Device) {
	callMeter, errorMeter := metrics.GetCallMeters("device_manager.get_device")
	callMeter.Mark(1)
	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: deviceAlias}
	getItemRequest := dynamodb.GetItemRequest{Key: key, TableName: this.DeviceTable}

	item, err := this.DynamoClient.GetItem(getItemRequest)
	if err == nil {
		return nil, &model.Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
	} else {
		errorMeter.Mark(1)
		return err, nil
	}
}

func (this *DefaultDeviceManager) GetDevices(cursor string) (error, *model.DeviceList) {
	callMeter, errorMeter := metrics.GetCallMeters("device_manager.get_devices")
	callMeter.Mark(1)
	var startKey map[string]dynamodb.Attribute
	if len(cursor) > 0 {
		startKey = make(map[string]dynamodb.Attribute)
		startKey["alias"] = dynamodb.Attribute{S: cursor}
	}

	scanRequest := dynamodb.ScanRequest{ExclusiveStartKey: startKey, TableName: this.DeviceTable, Limit: 1000}

	response, err := this.DynamoClient.ScanForItems(scanRequest)
	if err == nil {
		cursor = ""
		if len(response.LastEvaluatedKey) != 0 {
			cursor = response.LastEvaluatedKey["alias"].S
		}
		return nil, &model.DeviceList{convertToDevices(response.Items), cursor}
	} else {
		errorMeter.Mark(1)
		return err, nil
	}

}

func convertToDevices(items []map[string]dynamodb.Attribute) []model.Device {

	devices := make([]model.Device, 0, 0)
	for _, item := range items {
		device := model.Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
		devices = append(devices, device)
	}

	return devices
}

func formatTags(locale string, alias string, tags []string) string {

	tagString := alias
	tagString = tagString + "," + locale
	for _, tag := range tags {
		tagString = tagString + "," + tag
	}
	return tagString
}

package access

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/metrics"
	"github.com/gopns/gopns/model"
)

type DeviceManager interface {
	ListAppDevices(appId string, cursor string) (devices *[]model.Device, newCursor string, err error)
	RegisterDevice(device model.Device) (int, error)
	GetDevice(deviceId string) (*model.Device, error)
	PutDevice(device model.Device) error
}

type DefaultDeviceManager struct {
	SnsClient    sns.SNSClient
	DynamoClient dynamodb.DynamoClient
	DeviceTable  string
}

func NewDeviceManager(snsClient sns.SNSClient, dynamoClient dynamodb.DynamoClient, deviceTable string) DeviceManager {
	deviceManagerInstance := &DefaultDeviceManager{
		snsClient,
		dynamoClient,
		deviceTable}
	return deviceManagerInstance
}

func (dm *DefaultDeviceManager) ListAppDevices(appId string, cursor string) (devices *[]model.Device, newCursor string, err error) {
	return nil, "", nil
}

func (this *DefaultDeviceManager) RegisterDevice(device model.Device) (int, error) {
	callMeter, errorMeter := metrics.GetCallMeters("device_manager.register_device")
	callMeter.Mark(1)

	err := model.ValidateLocale(device.Locale())
	if err != nil {
		errorMeter.Mark(1)
		return 400, err
	}

	// TODO dummy app ARN replace by app manager
	arn, err := this.SnsClient.RegisterDevice(device.Token(), "", "DUMMY APP ARN")

	if err != nil {
		errorMeter.Mark(1)
		return 400, err
	}

	key := make(map[string]dynamodb.Attribute)
	key["id"] = dynamodb.Attribute{S: device.Id()}
	attributeUpdates := make(map[string]dynamodb.AttributeUpdate)
	attributeUpdates["arn"] = dynamodb.AttributeUpdate{"ADD", dynamodb.Attribute{SS: []string{arn}}}
	attributeUpdates["locale"] = dynamodb.AttributeUpdate{"PUT", dynamodb.Attribute{S: device.Locale()}}
	attributeUpdates["platform"] = dynamodb.AttributeUpdate{"PUT", dynamodb.Attribute{S: device.AppId()}}

	updateItemRequest := dynamodb.UpdateItemRequest{
		Key:              key,
		AttributeUpdates: attributeUpdates,
		TableName:        this.DeviceTable}
	err = this.DynamoClient.UpdateItem(updateItemRequest)

	if err != nil {
		errorMeter.Mark(1)
		return 500, err
	}

	return 0, nil
}

func (this *DefaultDeviceManager) GetDevice(deviceAlias string) (*model.Device, error) {
	callMeter, errorMeter := metrics.GetCallMeters("device_manager.get_device")
	callMeter.Mark(1)
	key := make(map[string]dynamodb.Attribute)
	// ToDo Fix
	key["alias"] = dynamodb.Attribute{S: deviceAlias}
	getItemRequest := dynamodb.GetItemRequest{Key: key, TableName: this.DeviceTable}

	item, err := this.DynamoClient.GetItem(getItemRequest)
	if err == nil {
		d := model.Device{}
		d.SetUserAlias(item["alias"].S)
		d.SetLocale(item["locale"].S)
		d.SetArn(item["arn"].S)
		d.SetDeviceType(model.DeviceType(item["deviceType"].S))
		return &d, nil
	} else {
		errorMeter.Mark(1)
		return nil, err
	}
}

func (dm *DefaultDeviceManager) PutDevice(device model.Device) error {
	return nil
}

/*
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
*/

func convertToDevices(items []map[string]dynamodb.Attribute) []model.Device {

	devices := make([]model.Device, 0, 0)
	for _, item := range items {
		d := model.Device{}
		d.SetUserAlias(item["alias"].S)
		d.SetLocale(item["locale"].S)
		d.SetArn(item["arn"].S)
		d.SetDeviceType(model.DeviceType(item["deviceType"].S))
		devices = append(devices, d)
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

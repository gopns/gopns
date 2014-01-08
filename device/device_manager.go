package device

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/gopnsconfig"
)

type DeviceManagerStruct struct {
	AwsConfig gopnsconfig.AWSConfig
}

type DeviceManager interface {
	RegisterDevice(device DeviceRegistration) (error, int)
	GetDevice(deviceAlias string) (error, *Device)
	GetDevices(cursor string) (error, *DeviceList)
}

var deviceManagerInstance DeviceManager

func InitilizeDeviceManager(config gopnsconfig.AWSConfig) {
	deviceManagerInstance = &DeviceManagerStruct{config}
}

func DeviceManagerInstance() DeviceManager {
	return deviceManagerInstance
}

func (this *DeviceManagerStruct) RegisterDevice(device DeviceRegistration) (error, int) {
	err := device.ValidateLocale()
	if err != nil {
		return err, 400
	}

	arn, err := sns.RegisterDevice(
		device.Id,
		formatTags(device.Locale, device.Alias, device.Tags),
		this.AwsConfig.UserID(),
		this.AwsConfig.UserSecret(),
		this.AwsConfig.Region(),
		this.AwsConfig.PlatformApps()[device.PlatformApp].Arn())

	if err != nil {
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
		TableName:        this.AwsConfig.DynamoTable()}
	err = dynamodb.UpdateItem(
		updateItemRequest,
		this.AwsConfig.UserID(),
		this.AwsConfig.UserSecret(),
		this.AwsConfig.Region())

	if err != nil {
		return err, 500
	}

	return nil, 0
}

func (this *DeviceManagerStruct) GetDevice(deviceAlias string) (error, *Device) {
	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: deviceAlias}
	getItemRequest := dynamodb.GetItemRequest{Key: key, TableName: this.AwsConfig.DynamoTable()}

	item, err := dynamodb.GetItem(
		getItemRequest,
		this.AwsConfig.UserID(),
		this.AwsConfig.UserSecret(),
		this.AwsConfig.Region())

	if err == nil {
		return nil, &Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
	} else {
		return err, nil
	}
}

func (this *DeviceManagerStruct) GetDevices(cursor string) (error, *DeviceList) {
	var startKey map[string]dynamodb.Attribute
	if len(cursor) > 0 {
		startKey = make(map[string]dynamodb.Attribute)
		startKey["alias"] = dynamodb.Attribute{S: cursor}
	}

	scanRequest := dynamodb.ScanRequest{ExclusiveStartKey: startKey, TableName: this.AwsConfig.DynamoTable(), Limit: 1000}

	response, err := dynamodb.ScanForItems(
		scanRequest,
		this.AwsConfig.UserID(),
		this.AwsConfig.UserSecret(),
		this.AwsConfig.Region())
	if err == nil {
		cursor = ""
		if len(response.LastEvaluatedKey) != 0 {
			cursor = response.LastEvaluatedKey["alias"].S
		}
		return nil, &DeviceList{convertToDevices(response.Items), cursor}
	} else {
		return err, nil
	}

}

func convertToDevices(items []map[string]dynamodb.Attribute) []Device {

	devices := make([]Device, 0, 0)
	for _, item := range items {
		device := Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
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

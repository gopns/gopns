package device

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/gopnsconfig"
)

type RegistrarStruct struct {
	AwsConfig gopnsconfig.AWSConfig
}

type Registrar interface {
	RegisterDevice(device DeviceRegistration) (error, int)
}

var registrarInstance Registrar

func InitilizeRegistrar(config gopnsconfig.AWSConfig) {
	registrarInstance = &RegistrarStruct{config}
}

func RegistrarInstance() Registrar {
	return registrarInstance
}

func (this *RegistrarStruct) RegisterDevice(device DeviceRegistration) (error, int) {
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

func formatTags(locale string, alias string, tags []string) string {

	tagString := alias
	tagString = tagString + "," + locale
	for _, tag := range tags {
		tagString = tagString + "," + tag
	}
	return tagString
}

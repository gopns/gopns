package gpnsconfig

import (
	"github.com/msbranco/goconfig"
)

var awsConfigInstance AWSConfig

func AWSConfigInstance() AWSConfig {
	return awsConfigInstance
}

type AWSConfigStruct struct {
	UserIDValue       string
	UserSecretValue   string
	PlatformAppsValue map[string]PlatformApp
	DynamoTableValue  string
	RegionValue       string
}

func (this AWSConfigStruct) UserID() string {
	return this.UserIDValue
}

func (this AWSConfigStruct) UserSecret() string {
	return this.UserSecretValue
}

func (this AWSConfigStruct) PlatformApps() map[string]PlatformApp {
	return this.PlatformAppsValue
}

func (this AWSConfigStruct) DynamoTable() string {
	return this.DynamoTableValue
}

func (this AWSConfigStruct) Region() string {
	return this.RegionValue
}

type AWSConfig interface {
	UserID() string
	UserSecret() string
	PlatformApps() map[string]PlatformApp
	DynamoTable() string
	Region() string
}

func parseAwsConfig(awsConfig *goconfig.ConfigFile) {
	userId, err := awsConfig.GetString("default", "id")
	checkError("Unable to find AWS User ID", err)

	userSecret, err := awsConfig.GetString("default", "secret")
	checkError("Unable to find AWS User Secret", err)

	dynamoTableValue, err := awsConfig.GetString("default", "dynamo-table")
	checkError("Unable to find AWS Dynamo Table", err)

	region, err := awsConfig.GetString("default", "region")
	checkError("Unable to find AWS region", err)

	awsConfigInstance = AWSConfigStruct{userId, userSecret, parsePlatformAppConfig(awsConfig), dynamoTableValue, region}
}

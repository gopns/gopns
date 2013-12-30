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

type AWSConfig interface {
	UserID() string
	UserSecret() string
	PlatformApps() map[string]PlatformApp
}

func parseAwsConfig(awsConfig *goconfig.ConfigFile) {
	userId, err := awsConfig.GetString("default", "id")
	checkError("Unable to find AWS User ID", err)

	userSecret, err := awsConfig.GetString("default", "secret")
	checkError("Unable to find AWS User Secret", err)

	awsConfigInstance = AWSConfigStruct{userId, userSecret, parsePlatformAppConfig(awsConfig)}
}

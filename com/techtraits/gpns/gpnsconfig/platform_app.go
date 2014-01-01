package gpnsconfig

import (
	"github.com/msbranco/goconfig"
	"strings"
)

type PlatformAppStruct struct {
	ArnValue  string
	TypeValue string
}

func (this PlatformAppStruct) Arn() string {
	return this.ArnValue
}

func (this PlatformAppStruct) Type() string {
	return this.TypeValue
}

type PlatformApp interface {
	Arn() string
	Type() string
}

func parsePlatformAppConfig(awsConfig *goconfig.ConfigFile) map[string]PlatformApp {
	platformApps, err := awsConfig.GetString("default", "platform-applications")
	checkError("Unable to find AWS Platform Apps List", err)
	platformAppsMap := make(map[string]PlatformApp)
	for _, platformApp := range strings.Split(platformApps, ",") {
		arn, err := awsConfig.GetString(platformApp, "arn")
		checkError("Unable to find AWS ARN for app "+platformApp, err)

		typeValue, err := awsConfig.GetString(platformApp, "type")
		checkError("Unable to find AWS type for app "+platformApp, err)

		platformAppsMap[platformApp] = PlatformAppStruct{arn, typeValue}

	}
	return platformAppsMap
}

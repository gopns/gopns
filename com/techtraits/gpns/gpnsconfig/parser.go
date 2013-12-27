package gpnsconfig

import (
	"flag"
	"github.com/msbranco/goconfig"
	"log"
	"strings"
)

func ParseConfig() (BaseConfig, AWSConfig) {
	var aws_config_file string
	var base_config_file string

	flag.StringVar(&base_config_file, "baseConfig", "./config/base.conf", "The path to the base configuration file")
	flag.StringVar(&aws_config_file, "awsConfig", "./config/aws.conf", "The path to the aws configuration file")
	flag.Parse()

	log.Printf("Using base configuration file: %s", base_config_file)
	baseConfig, err := goconfig.ReadConfigFile(base_config_file)
	checkError("Unable to parse base config", err)

	log.Printf("Using aws configuration file: %s", aws_config_file)
	awsConfig, err := goconfig.ReadConfigFile(aws_config_file)
	checkError("Unable to parse AWS config", err)

	port, err := baseConfig.GetString("default", "port")
	checkError("Unable to find Server Port", err)

	userId, err := awsConfig.GetString("default", "id")
	checkError("Unable to find AWS User ID", err)

	userSecret, err := awsConfig.GetString("default", "secret")
	checkError("Unable to find AWS User Secret", err)

	//region, err := awsConfig.GetString("default", "region")
	//checkError("Unable to find AWS Region", err)

	platformApps, err := awsConfig.GetString("default", "platform-applications")
	checkError("Unable to find AWS Platform Apps List", err)

	platformAppsMap := make(map[string]PlatformApp)
	for _, platformApp := range strings.Split(platformApps, ",") {
		arn, err := awsConfig.GetString(platformApp, "arn")
		checkError("Unable to find AWS ARN for app "+platformApp, err)
		region, err := awsConfig.GetString(platformApp, "region")
		checkError("Unable to find AWS region for app "+platformApp, err)
		platformAppsMap[platformApp] = PlatformAppStruct{arn, region}

	}

	return BaseConfigStruct{port}, AWSConfigStruct{userId, userSecret, platformAppsMap}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

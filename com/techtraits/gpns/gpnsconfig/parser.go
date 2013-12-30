package gpnsconfig

import (
	"flag"
	"github.com/msbranco/goconfig"
	"log"
)

func ParseConfig() {
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

	parseBaseConfig(baseConfig)
	parseAwsConfig(awsConfig)

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

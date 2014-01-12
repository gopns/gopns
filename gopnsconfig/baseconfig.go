package gopnsconfig

import (
	"log"
)

type BaseConfigStruct struct {
	PortValue          string
	MetricsServerValue string
	MetricsAPIKeyValue string
	MetricsPrefixValue string
}

func (this *BaseConfigStruct) Port() string {
	return this.PortValue
}

func (this *BaseConfigStruct) MetricsAPIKey() string {
	return this.MetricsAPIKeyValue
}

func (this *BaseConfigStruct) MetricsServer() string {
	return this.MetricsServerValue
}

func (this *BaseConfigStruct) MetricsPrefix() string {
	return this.MetricsPrefixValue
}

type BaseConfig interface {
	Port() string
	MetricsServer() string
	MetricsAPIKey() string
	MetricsPrefix() string
}

func parseBaseConfig(baseConfig *ConfigFile) BaseConfig {
	port, err := baseConfig.GetString("default", "port")
	checkError("Unable to find Server Port", err)

	metricsServer, err := baseConfig.GetString("default", "metrics-server")
	checkError("Unable to find metrics server", err)

	metricsKey, err := baseConfig.GetString("default", "metrics-api-key")
	if err != nil {
		log.Println("Unable to find metrics-api-key using empty string")
		metricsKey = ""
	}

	metricsPrefix, err := baseConfig.GetString("default", "metrics-prefix")
	if err != nil {
		log.Println("Unable to find metrics prefix using empty string")
		metricsPrefix = ""
	}

	baseConfigInstance := &BaseConfigStruct{port, metricsServer, metricsKey, metricsPrefix}
	return baseConfigInstance

}

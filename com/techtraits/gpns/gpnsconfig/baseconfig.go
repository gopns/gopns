package gpnsconfig

import (
	"github.com/msbranco/goconfig"
)

var baseConfigInstance BaseConfig

func BaseConfigInstance() BaseConfig {
	return baseConfigInstance
}

type BaseConfigStruct struct {
	PortValue string
}

func (this BaseConfigStruct) Port() string {
	return this.PortValue
}

type BaseConfig interface {
	Port() string
}

func parseBaseConfig(baseConfig *goconfig.ConfigFile) {
	port, err := baseConfig.GetString("default", "port")
	checkError("Unable to find Server Port", err)
	baseConfigInstance = BaseConfigStruct{port}
}

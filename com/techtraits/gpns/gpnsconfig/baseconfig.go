package gpnsconfig

import (
	"github.com/msbranco/goconfig"
)

type BaseConfigStruct struct {
	PortValue string
}

func (this BaseConfigStruct) Port() string {
	return this.PortValue
}

type BaseConfig interface {
	Port() string
}

func parseBaseConfig(baseConfig *goconfig.ConfigFile) BaseConfigStruct {
	port, err := baseConfig.GetString("default", "port")
	checkError("Unable to find Server Port", err)
	return BaseConfigStruct{port}
}

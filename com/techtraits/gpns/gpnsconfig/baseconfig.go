package gpnsconfig

type BaseConfigStruct struct {
	PortValue string
}

func (this BaseConfigStruct) Port() string {
	return this.PortValue
}

type BaseConfig interface {
	Port() string
}

package modelview

import "github.com/gopns/gopns/model"

type DeviceRegistration struct {
	userId      string
	Id          string
	Locale      string
	PlatformApp string
	Tags        []string
}

type DeviceList struct {
	Devices []model.Device
	Cursor  string `json:",omitempty"`
}

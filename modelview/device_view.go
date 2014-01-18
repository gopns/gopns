package modelview

import (
	"errors"
	"github.com/gopns/gopns/model"
	"regexp"
)

type RegisterDeviceView struct {
	userAlias  string
	deviceType model.DeviceType
	token      string
	locale     string
	timezone   string
}

func NewRegisterDeviceView(userAlias string, dt model.DeviceType, token string, locale string, timezone string) *NewRegisterDeviceView {
	return &NewRegisterDeviceView{userAlias: userAlias, token: token, locale: locale, timezone: timezone}
}

func (dv RegisterDeviceView) UserAlias() string {
	return dv.userAlias
}

func (dv RegisterDeviceView) Locale() string {
	return dv.locale
}

func (dv RegisterDeviceView) DeviceType() model.DeviceType {
	return dv.deviceType
}

func (dv RegisterDeviceView) Token() string {
	return dv.token
}

func (dv RegisterDeviceView) Timezone() string {
	return dv.timezone
}

func (dv *RegisterDeviceView) SetUserAlias(id string) {
	dv.userAlias = id
}

func (dv *RegisterDeviceView) SetLocale(locale string) {
	dv.locale = locale
}

func (dv *RegisterDeviceView) SetTimezone(timezone string) {
	dv.timezone = timezone
}

func (dv *RegisterDeviceView) SetDeviceType(dt model.DeviceType) {
	dv.deviceType = dt
}

func (dv *RegisterDeviceView) SetToken(token string) {
	dv.token = token
}

type DeviceList struct {
	Devices []model.Device
	Cursor  string `json:",omitempty"`
}

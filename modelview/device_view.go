package modelview

import (
	"github.com/gopns/gopns/model"
)

//Eventually they will diverge
type DeviceView model.Device

func FromDevice(d model.Device) *DeviceView {
	dv := DeviceView(d)
	return &dv
}

type DeviceRegisterView struct {
	userAlias  string
	deviceType model.DeviceType
	token      string
	locale     string
	timezone   string
}

func NewDeviceRegisterView(userAlias string, dt model.DeviceType, token string, locale string, timezone string) *DeviceRegisterView {
	return &DeviceRegisterView{userAlias: userAlias, token: token, locale: locale, timezone: timezone}
}

func (dv DeviceRegisterView) UserAlias() string {
	return dv.userAlias
}

func (dv DeviceRegisterView) Locale() string {
	return dv.locale
}

func (dv DeviceRegisterView) DeviceType() model.DeviceType {
	return dv.deviceType
}

func (dv DeviceRegisterView) Token() string {
	return dv.token
}

func (dv DeviceRegisterView) Timezone() string {
	return dv.timezone
}

func (dv *DeviceRegisterView) SetUserAlias(id string) {
	dv.userAlias = id
}

func (dv *DeviceRegisterView) SetLocale(locale string) {
	dv.locale = locale
}

func (dv *DeviceRegisterView) SetTimezone(timezone string) {
	dv.timezone = timezone
}

func (dv *DeviceRegisterView) SetDeviceType(dt model.DeviceType) {
	dv.deviceType = dt
}

func (dv *DeviceRegisterView) SetToken(token string) {
	dv.token = token
}

func (dv *DeviceRegisterView) ToDevice() (*model.Device, error) {
	device := new(model.Device)
	device.SetUserAlias(dv.userAlias)
	device.SetToken(dv.token)
	device.SetTimezone(dv.timezone)
	if err := device.SetDeviceType(dv.deviceType); err != nil {
		return nil, err
	}
	if err := device.SetLocale(dv.locale); err != nil {
		return nil, err
	}

	return device, nil
}

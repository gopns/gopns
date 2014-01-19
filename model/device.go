package model

import (
	"errors"
	"regexp"
)

type DeviceType string

const (
	IOS     DeviceType = "IOS"
	ANDROID DeviceType = "ANDROID"
	KINDLE  DeviceType = "KINDLE"
)

//TODO: constants for locale and timezone
type Device struct {
	id         string
	userAlias  string
	appId      string
	deviceType DeviceType
	arn        string
	token      string
	locale     string
	timezone   string
	enabled    bool
}

func NewDevice(id string, userAlias string, appId string, dt DeviceType, arn string, token string, locale string, timezone string, enabled bool) (*Device, error) {
	device := &Device{id: id, userAlias: userAlias, appId: appId, arn: arn, token: token, enabled: enabled}
	if err := device.SetDeviceType(dt); err != nil {
		return nil, err
	}
	if err := device.SetLocale(locale); err != nil {
		return nil, err
	}
	//TODO validate
	device.SetTimezone(timezone)
	return device, nil
}

func (device Device) Id() string {
	return device.id
}

func (device Device) UserAlias() string {
	return device.userAlias
}

func (device Device) AppId() string {
	return device.appId
}

func (device Device) Locale() string {
	return device.locale
}

func (device Device) DeviceType() DeviceType {
	return device.deviceType
}

func (device Device) Arn() string {
	return device.arn
}

func (device Device) Token() string {
	return device.token
}

func (device Device) Enabled() bool {
	return device.enabled
}

func (device Device) Timezone() string {
	return device.timezone
}

func (device Device) Platform() Platform {
	if device.deviceType == IOS {
		return APNS
	} else if device.deviceType == ANDROID {
		return GCM
	} else if device.deviceType == KINDLE {
		return ADM
	}
	return UNKNOWN
}

func (device *Device) SetId(id string) {
	device.id = id
}

func (device *Device) SetUserAlias(id string) {
	device.userAlias = id
}

func (device *Device) SetAppId(id string) {
	device.appId = id
}

func (device *Device) SetLocale(locale string) error {
	if err := ValidateLocale(locale); err != nil {
		return err
	}
	device.locale = locale
	return nil
}

func (device *Device) SetTimezone(timezone string) {
	device.timezone = timezone
}

func (device *Device) SetDeviceType(dt DeviceType) error {
	if err := ValidateDeviceType(dt); err != nil {
		return err
	}
	device.deviceType = dt
	return nil
}

func (device *Device) SetArn(arn string) {
	device.arn = arn
}

func (device *Device) SetToken(token string) {
	device.token = token
}

func (device *Device) SetEnabled(enabled bool) {
	device.enabled = enabled
}

var localRegex = regexp.MustCompile(`^[A-Za-z]{2,3}_[A-Za-z]{2,3}$`)

func ValidateLocale(locale string) error {
	if !localRegex.MatchString(locale) {
		return errors.New("Invalid locale, should be of the form US_en.")
	}

	return nil
}

var deviceTypeRegex = regexp.MustCompile("^(" + string(ANDROID) + "|" + string(IOS) + "|" + string(KINDLE) + ")$")

func ValidateDeviceType(dt DeviceType) error {
	if !deviceTypeRegex.MatchString(string(dt)) {
		return errors.New("Invalid device type, valid values: " + string(ANDROID) + "," + string(IOS) + " or " + string(KINDLE))
	}

	return nil
}

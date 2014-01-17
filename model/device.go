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
	userId     string
	appId      string
	deviceType DeviceType
	arn        string
	token      string
	locale     string
	timezone   string
	enabled    bool
}

func (device Device) UserId() string {
	return device.userId
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

func (device *Device) SetUserId(id string) {
	device.userId = id
}

func (device *Device) SetAppId(id string) {
	device.appId = id
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

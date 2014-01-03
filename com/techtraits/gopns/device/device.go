package device

import (
	"errors"
	"regexp"
)

var validLocale = regexp.MustCompile(`^[A-Za-z]{2,3}_[A-Za-z]{2,3}$`)

type DeviceRegistration struct {
	Alias       string
	Id          string
	Locale      string
	PlatformApp string
	Tags        []string
}

type Device struct {
	Alias  string
	Locale string
	Arns   []string
	Tags   []string
}

type DeviceList struct {
	Devices []Device
	cursor  string
}

func (this DeviceRegistration) ValidateLocale() error {
	//TODO More validation to actually check codes
	if !validLocale.MatchString(this.Locale) {
		return errors.New("Invalid locale")
	}

	return nil

}

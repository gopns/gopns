package device

import (
	"errors"
	"regexp"
)

var validLocale = regexp.MustCompile(`^[A-Za-z]{2,3}_[A-Za-z]{2,3}$`)

//TODO Verify IOS and ADM as keys
var validPlatform = regexp.MustCompile(`^(GCM|ADM|IOS)$`)

type DeviceRegistration struct {
	Alias       string
	Id          string
	Locale      string
	PlatformApp string
	Tags        []string
}

type Device struct {
	Alias    string
	Locale   string
	Arns     []string
	Platform string
	Tags     []string
}

type DeviceList struct {
	Devices []Device
	Cursor  string
}

func (this DeviceRegistration) ValidateLocale() error {
	//TODO More validation to actually check codes
	if !validLocale.MatchString(this.Locale) {
		return errors.New("Invalid locale")
	}

	return nil

}

func ValidateLocale(locale string) error {
	//TODO More validation to actually check codes
	if !validLocale.MatchString(locale) {
		return errors.New("Invalid locale")
	}

	return nil
}

func ValidatePlatform(platform string) error {
	//TODO More validation to actually check codes
	if !validPlatform.MatchString(platform) {
		return errors.New("Invalid platform")
	}

	return nil
}

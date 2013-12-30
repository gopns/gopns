package device

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

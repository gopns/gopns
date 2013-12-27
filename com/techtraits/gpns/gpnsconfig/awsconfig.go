package gpnsconfig

type AWSConfigStruct struct {
	UserIDValue       string
	UserSecretValue   string
	PlatformAppsValue map[string]PlatformApp
}

func (this AWSConfigStruct) UserID() string {
	return this.UserIDValue
}

func (this AWSConfigStruct) UserSecret() string {
	return this.UserSecretValue
}

func (this AWSConfigStruct) PlatformApps() map[string]PlatformApp {
	return this.PlatformAppsValue
}

type AWSConfig interface {
	UserID() string
	UserSecret() string
	PlatformApps() map[string]PlatformApp
}

type PlatformAppStruct struct {
	ArnValue    string
	RegionValue string
}

func (this PlatformAppStruct) Arn() string {
	return this.ArnValue
}

func (this PlatformAppStruct) Region() string {
	return this.RegionValue
}

type PlatformApp interface {
	Arn() string
	Region() string
}

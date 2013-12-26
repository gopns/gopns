package gpnsconfig

type AWSConfigStruct struct {
	UserIDValue     string
	UserSecretValue string
}

func (this AWSConfigStruct) UserID() string {
	return this.UserIDValue
}

func (this AWSConfigStruct) UserSecret() string {
	return this.UserSecretValue
}

type AWSConfig interface {
	UserID() string
	UserSecret() string
}

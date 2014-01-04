package gopnsconfig

import (
	"github.com/msbranco/goconfig"
)

var awsConfigInstance AWSConfig

func AWSConfigInstance() AWSConfig {
	return awsConfigInstance
}

type AWSConfigStruct struct {
	UserIDValue               string
	UserSecretValue           string
	PlatformAppsValue         map[string]PlatformApp
	DynamoTableValue          string
	RegionValue               string
	InitialReadCapacityValue  int
	InitialWriteCapacityValue int
	SqsQueueNameValue         string
	SqsQueueUrlValue          string
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

func (this AWSConfigStruct) DynamoTable() string {
	return this.DynamoTableValue
}

func (this AWSConfigStruct) Region() string {
	return this.RegionValue
}

func (this AWSConfigStruct) InitialReadCapacity() int {
	return this.InitialReadCapacityValue
}

func (this AWSConfigStruct) InitialWriteCapacity() int {
	return this.InitialWriteCapacityValue
}

func (this AWSConfigStruct) SqsQueueName() string {
	return this.SqsQueueNameValue
}

func (this AWSConfigStruct) SqsQueueUrl() string {
	return this.SqsQueueUrlValue
}

func (this AWSConfigStruct) SetSqsQueueUrl(queueUrl string) {
	this.SqsQueueUrlValue = queueUrl
}

type AWSConfig interface {
	UserID() string
	UserSecret() string
	PlatformApps() map[string]PlatformApp
	DynamoTable() string
	Region() string
	InitialReadCapacity() int
	InitialWriteCapacity() int
	SqsQueueName() string
	SqsQueueUrl() string
	SetSqsQueueUrl(queueUrl string)
}

func parseAwsConfig(awsConfig *goconfig.ConfigFile) {
	userId, err := awsConfig.GetString("default", "id")
	checkError("Unable to find AWS User ID", err)

	userSecret, err := awsConfig.GetString("default", "secret")
	checkError("Unable to find AWS User Secret", err)

	dynamoTableValue, err := awsConfig.GetString("default", "dynamo-table")
	checkError("Unable to find AWS Dynamo Table", err)

	region, err := awsConfig.GetString("default", "region")
	checkError("Unable to find AWS region", err)

	readCapacity, err := awsConfig.GetInt64("default", "dynamo-read-capacity")
	checkError("Unable to find AWS dynamo-read-capacity", err)

	writeCapacity, err := awsConfig.GetInt64("default", "dynamo-write-capacity")
	checkError("Unable to find AWS dynamo-write-capacity", err)

	sqsQueueName, err := awsConfig.GetString("default", "sqs-queue-name")
	checkError("Unable to find AWS sqs-queue-name", err)

	awsConfigInstance = AWSConfigStruct{
		userId,
		userSecret,
		parsePlatformAppConfig(awsConfig),
		dynamoTableValue,
		region,
		int(readCapacity),
		int(writeCapacity),
		sqsQueueName,
		""}
}

package gopnsconfig

var awsConfigInstance AWSConfig

type AWSConfigStruct struct {
	awsAccessKeyId               string
	awsSecretKey           		 string
	region					 string
	PlatformAppsValue         map[string]PlatformApp
	DynamoTableValue          string
	InitialReadCapacityValue  int
	InitialWriteCapacityValue int
	SqsQueueNameValue         string
	SqsQueueUrlValue          string
}

func (this *AWSConfigStruct) AwsAccessKeyId() string {
	return this.awsAccessKeyId
}

func (this *AWSConfigStruct) AwsSecretKey() string {
	return this.awsSecretKey
}

func (this *AWSConfigStruct) Region() string {
	return this.region
}

func (this *AWSConfigStruct) PlatformApps() map[string]PlatformApp {
	return this.PlatformAppsValue
}

func (this *AWSConfigStruct) PlatformAppsMap() map[string]map[string]string {
	platformAppsMap := make(map[string]map[string]string)
	for appName, app := range this.PlatformAppsValue {
		platformAppsMap[appName] = app.ConfigMap()
	}
	return platformAppsMap
}

func (this *AWSConfigStruct) DynamoTable() string {
	return this.DynamoTableValue
}



func (this *AWSConfigStruct) InitialReadCapacity() int {
	return this.InitialReadCapacityValue
}

func (this *AWSConfigStruct) InitialWriteCapacity() int {
	return this.InitialWriteCapacityValue
}

func (this *AWSConfigStruct) SqsQueueName() string {
	return this.SqsQueueNameValue
}

func (this *AWSConfigStruct) SqsQueueUrl() string {
	return this.SqsQueueUrlValue
}

func (this *AWSConfigStruct) SetSqsQueueUrl(queueUrl string) {
	(*this).SqsQueueUrlValue = queueUrl
}

type AWSConfig interface {
	AwsAccessKeyId() string
	AwsSecretKey() string
	Region() string
	PlatformApps() map[string]PlatformApp
	PlatformAppsMap() map[string]map[string]string
	DynamoTable() string
	InitialReadCapacity() int
	InitialWriteCapacity() int
	SqsQueueName() string
	SqsQueueUrl() string
	SetSqsQueueUrl(queueUrl string)
}

func parseAwsConfig(awsConfig *ConfigFile) AWSConfig {
	awsAccessKeyId, err := awsConfig.GetString("default", "aws_access_key_id")
	checkError("Unable to find AWS Access Key ID (aws_access_key_id)", err)

	awsSecretKey, err := awsConfig.GetString("default", "aws_access_key_secret")
	checkError("Unable to find AWS access key secret (aws_access_key_secret", err)

	region, err := awsConfig.GetString("default", "region")
	checkError("Unable to find AWS region", err)

	dynamoTableValue, err := awsConfig.GetString("default", "dynamo-table")
	checkError("Unable to find AWS Dynamo Table", err)

	readCapacity, err := awsConfig.GetInt64("default", "dynamo-read-capacity")
	checkError("Unable to find AWS dynamo-read-capacity", err)

	writeCapacity, err := awsConfig.GetInt64("default", "dynamo-write-capacity")
	checkError("Unable to find AWS dynamo-write-capacity", err)

	sqsQueueName, err := awsConfig.GetString("default", "sqs-queue-name")
	checkError("Unable to find AWS sqs-queue-name", err)

	awsConfigInstance := &AWSConfigStruct{
		awsAccessKeyId,
		awsSecretKey,
		region,		
		parsePlatformAppConfig(awsConfig),
		dynamoTableValue,
		int(readCapacity),
		int(writeCapacity),
		sqsQueueName,
		""}
	return awsConfigInstance

}

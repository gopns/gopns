package main

import (
	"github.com/gopns/gopns/com/techtraits/gopns/aws/dynamodb"
	"github.com/gopns/gopns/com/techtraits/gopns/aws/sqs"
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/rest"
	"log"
)

func main() {

	appMode := gopnsconfig.ParseConfig()
	err := dynamodb.Initilize(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().DynamoTable(),
		gopnsconfig.AWSConfigInstance().InitialReadCapacity(),
		gopnsconfig.AWSConfigInstance().InitialWriteCapacity())

	if err != nil {
		log.Fatalf("Unable to initilize Dynamo DB %s", err.Error())
	}

	err, sqsQueue := sqs.Initilize(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().SqsQueueName())

	if err != nil {
		log.Fatalf("Unable to initilize SQS %s", err.Error())
	} else {
		log.Printf("Using SQS Queue %s", sqsQueue.QueueUrl)
		gopnsconfig.AWSConfigInstance().SetSqsQueueUrl(sqsQueue.QueueUrl)
	}

	//TODO Just here for testing delete
	err, _ = sqs.SendMessage(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().SqsQueueUrl(), "Test Message")

	sqs.SendMessage(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().SqsQueueUrl(), "Test Message 1")

	sqs.DeleteMessages(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().SqsQueueUrl(), sqsMessages)

	//End TODO

	if appMode == gopnsconfig.SERVER_MODE {
		rest.SetupRestServices()
	} else if appMode == gopnsconfig.REGISTER_MODE {

	} else if appMode == gopnsconfig.SEND_MODE {

	}

}

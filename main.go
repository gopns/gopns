package main

import (
	"github.com/gopns/gopns/gopnsapp"
	"github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/rest"
	"log"
)

func main() {

	appMode := gopnsconfig.ParseConfig()

	/*


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

		device.InitilizeDeviceManager(gopnsconfig.AWSConfigInstance())

		//TODO Just here for testing delete
		err, _ = sqs.SendMessage(gopnsconfig.AWSConfigInstance().UserID(),
			gopnsconfig.AWSConfigInstance().UserSecret(),
			gopnsconfig.AWSConfigInstance().Region(),
			gopnsconfig.AWSConfigInstance().SqsQueueUrl(), "Test Message")

		sqs.SendMessage(gopnsconfig.AWSConfigInstance().UserID(),
			gopnsconfig.AWSConfigInstance().UserSecret(),
			gopnsconfig.AWSConfigInstance().Region(),
			gopnsconfig.AWSConfigInstance().SqsQueueUrl(), "Test Message 1")

	*/
	/*
		sqs.DeleteMessages(gopnsconfig.AWSConfigInstance().UserID(),
			gopnsconfig.AWSConfigInstance().UserSecret(),
			gopnsconfig.AWSConfigInstance().Region(),
			gopnsconfig.AWSConfigInstance().SqsQueueUrl(), sqsMessages)
	*/
	//END TODO

	// start gopns -- MUST BE STARTED BEFORE ANYTHING ELSE
	gopnsapp_, err := gopnsapp.Initilize()
	if err == nil {
		gopnsapp_.Start()

		if appMode == gopnsconfig.SERVER_MODE {
			rest.SetupRestServices()
		} else if appMode == gopnsconfig.REGISTER_MODE {

		} else if appMode == gopnsconfig.SEND_MODE {

		}
	} else {
		log.Fatal(err.Error())
	}

}

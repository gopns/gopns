package notification

import (
	"encoding/json"
	"github.com/gopns/gopns/aws/sqs"
	"sync"
)

type NotificationConsumer interface {
	Start()
	Stop()
}

const (
	stopCMD  = iota - 1 // stopCMD == -1
	startCMD = iota     // startCMD == 0
)

type SQSNotificationConsumer struct {
	sqsClient         sqs.SQSClient
	Sender            *NotificationSender
	sqsQueueUrl       string
	processorKillChan chan bool
	processor_wg      sync.WaitGroup
}

func NewSQSNotifictionConsumer(queueUrl string, sqsClient sqs.SQSClient, sender *NotificationSender) (consumer *SQSNotificationConsumer) {
	consumer = new(SQSNotificationConsumer)
	consumer.sqsQueueUrl = queueUrl
	consumer.sqsClient = sqsClient
	consumer.Sender = sender
	consumer.processorKillChan = make(chan bool, 1)
	return consumer
}

func (consumer *SQSNotificationConsumer) Start() {
	// ToDo check if the queue processor is already running or not
	consumer.processor_wg.Add(1)
	go consumer.processor()
}

func (consumer *SQSNotificationConsumer) Stop() {
	consumer.processorKillChan <- true
	consumer.processor_wg.Wait()
}

func (consumer *SQSNotificationConsumer) processor() {

	var sqsMessages []sqs.SqsMessage
QUEUE_PROCESS_LOOP:
	for {
		select {
		// stopping
		case <-consumer.processorKillChan:
			break QUEUE_PROCESS_LOOP
		default:
			// consume notificationt tasks from sqs client and use notification sender to distribute work

			sqsMessages, _ = consumer.sqsClient.GetMessage(consumer.sqsQueueUrl, 10, 20) //long polling, wait for upto 20 seconds before giving up

			var task NotificationTask
			for _, sqsMessage := range sqsMessages {
				_ = json.Unmarshal([]byte(sqsMessage.Body), &task)
				// TODO enable sending notifications after adding error handling
				//consumer.Sender.SendNotification(task)
			}

			// delete processed notification tasks from queue
			consumer.sqsClient.DeleteMessages(consumer.sqsQueueUrl, sqsMessages)

		}
	}
	consumer.processor_wg.Done()
}

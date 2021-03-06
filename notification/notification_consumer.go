package notification

import (
	"encoding/json"
	"github.com/gopns/gopns/aws/sqs"
	"github.com/gopns/gopns/metrics"
	"sync"
)

type NotificationConsumer interface {
	Start() error
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

func (this *SQSNotificationConsumer) Start() error {
	// ToDo check if the queue processor is already running or not
	this.processor_wg.Add(1)
	go this.processor()
	return nil
}

func (this *SQSNotificationConsumer) Stop() {
	this.processorKillChan <- true
	this.processor_wg.Wait()
}

func (this *SQSNotificationConsumer) processor() {

	var sqsMessages []sqs.SqsMessage
QUEUE_PROCESS_LOOP:
	for {
		select {
		// stopping
		case <-this.processorKillChan:
			break QUEUE_PROCESS_LOOP
		default:

			// consume notificationt tasks from sqs client and use notification sender to distribute work
			sqsMessages, _ = this.sqsClient.GetMessage(this.sqsQueueUrl, 10, 20) //long polling, wait for upto 20 seconds before giving up
			//TODO CHECK ERROR
			var task NotificationTask
			for _, sqsMessage := range sqsMessages {
				callMeter, _ := metrics.GetCallMeters("notification_consumer.message_consumed")
				callMeter.Mark(1)
				_ = json.Unmarshal([]byte(sqsMessage.Body), &task)
				// TODO enable sending notifications after adding error handling
				//this.Sender.SendNotification(task)
			}

			// delete processed notification tasks from queue
			this.sqsClient.DeleteMessages(this.sqsQueueUrl, sqsMessages)

		}
	}
	this.processor_wg.Done()
}

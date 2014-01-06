package notification

import (
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/stefantalpalaru/pool"
)

type NotificationSender struct {
	awsConfig  *AWSConfig
	workerPool *Pool
}

func (sender *NotificationSender) SendSyncNotification(device Device, message NotificationMessage, int timeout) int {
	c := &make(chan int, 1)

	task := NotificationTask{device: device, message: message, respondTo: c}
	SendNotification(this, task)
	select {
	case status := <-c:
		return status
	case <-time.After(timeout * time.Second):
		return 408 //timeout
	}

}

func (sender *NotificationSender) SendAsyncNotification(device Device, message NotificationMessage) {
	task := NotificationTask{device: device, message: message, respondTo: nil}
	SendNotification(this, task)
}

func (sender *NotificationSender) SendNotification(task NotificationTask) {
	workerPool.Add(sendNotification, this, task)
}

//actual function for sending notifications

func (sender *NotificationSender) sendNotification(task NotificationTask) {

	device = task.device
	message = task.message
	for _, arn := range device_.Arns {

		sns.PublishNotification(
			arn,
			message.Title,
			message.Message,
			awsConfig.UserID(),
			awsConfig.UserSecret(),
			awsConfig.Region(),
			awsConfig.PlatformApps()[device_.Platform].Type())
	}

	// send appropriate response code
	if c := task.respondTo; c != nil {
		c <- 202
	}

	return
}

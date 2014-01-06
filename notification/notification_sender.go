package notification

import (
	"github.com/gopns/aws/sns"
	"github.com/gopns/device"
	"github.com/gopns/gopnsconfig"
	"github.com/stefantalpalaru/pool"
	"time"
)

type NotificationSender struct {
	AwsConfig  gopnsconfig.AWSConfig
	WorkerPool *pool.Pool
}

func (sender *NotificationSender) SendSyncNotification(device device.Device, message NotificationMessage, timeout int) int {
	c := make(chan int, 1)

	task := NotificationTask{device: device, message: message, respondTo: &c}
	sender.SendNotification(task)
	select {
	case status := <-c:
		return status
	case <-time.After(time.Duration(timeout) * time.Second):
		return 408 //timeout
	}

}

func (sender *NotificationSender) SendAsyncNotification(device device.Device, message NotificationMessage) {
	task := NotificationTask{device: device, message: message, respondTo: nil}
	sender.SendNotification(task)
}

func (sender *NotificationSender) SendNotification(task NotificationTask) {
	sender.WorkerPool.Add(notificationWork, sender, task)
}

//actual function for sending notifications

func (sender *NotificationSender) sendNotification(task NotificationTask) {

	device_ := task.device
	message := task.message
	for _, arn := range device_.Arns {

		sns.PublishNotification(
			arn,
			message.Title,
			message.Message,
			sender.AwsConfig.UserID(),
			sender.AwsConfig.UserSecret(),
			sender.AwsConfig.Region(),
			sender.AwsConfig.PlatformApps()[device_.Platform].Type())
	}

	// send appropriate response code
	if c := *task.respondTo; c != nil {
		c <- 202
	}

	return
}

// a function closure for passing the job to a worker
func notificationWork(args ...interface{}) interface{} {
	sender := args[0].(*NotificationSender)
	task := args[1].(NotificationTask)
	sender.sendNotification(task)
	return nil
}

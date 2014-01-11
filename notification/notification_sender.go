package notification

import (
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/device"
	"github.com/gopns/gopns/gopnsconfig"
	"github.com/stefantalpalaru/pool"
	"time"
)

type NotificationSender struct {
	SnsClient    sns.SNSClient
	WorkerPool   *pool.Pool
	PlatformApps map[string]gopnsconfig.PlatformApp
}

func (this *NotificationSender) SendSyncNotification(device device.Device, message NotificationMessage, timeout int) int {
	c := make(chan int, 1)

	task := NotificationTask{device: device, message: message, respondTo: &c}
	this.SendNotification(task)
	select {
	case status := <-c:
		return status
	case <-time.After(time.Duration(timeout) * time.Second):
		return 408 //timeout
	}

}

func (this *NotificationSender) SendAsyncNotification(device device.Device, message NotificationMessage) {
	task := NotificationTask{device: device, message: message, respondTo: nil}
	this.SendNotification(task)
}

func (this *NotificationSender) SendNotification(task NotificationTask) {
	this.WorkerPool.Add(notificationWork, this, task)
}

//actual function for sending notifications

func (this *NotificationSender) sendNotification(task NotificationTask) {

	device_ := task.device
	message := task.message
	for _, arn := range device_.Arns {

		this.SnsClient.PublishNotification(
			arn,
			message.Title,
			message.Message,
			this.PlatformApps[device_.Platform].Type())
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
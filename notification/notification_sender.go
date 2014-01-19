package notification

import (
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/metrics"
	"github.com/gopns/gopns/model"
	"github.com/stefantalpalaru/pool"
	"time"
)

type NotificationSender struct {
	SnsClient    sns.SNSClient
	WorkerPool   *pool.Pool
	PlatformApps map[string]map[string]string
}

func (this *NotificationSender) SendSyncNotification(device model.Device, message NotificationMessage, timeout int) int {
	callMeter, errorMeter := metrics.GetCallMeters("notification_sender.send_sync_notification")
	callMeter.Mark(1)

	c := make(chan int, 1)

	task := NotificationTask{device: device, message: message, respondTo: &c}
	this.SendNotification(task)
	select {
	case status := <-c:
		return status
	case <-time.After(time.Duration(timeout) * time.Second):
		errorMeter.Mark(1)
		return 408 //timeout
	}

}

func (this *NotificationSender) SendAsyncNotification(device model.Device, message NotificationMessage) {
	callMeter, _ := metrics.GetCallMeters("notification_sender.send_async_notification")
	callMeter.Mark(1)
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
	callMeter, _ := metrics.GetCallMeters("notification_sender.send_notification")
	callMeter.Mark(1)
	this.SnsClient.PublishNotification(
		device_.Arn(),
		message.Title,
		message.Message,
		string(device_.Platform()))

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

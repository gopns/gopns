package notification

import (
	"github.com/gopns/gopns/device"
)

type NotificationTask struct {
	message   NotificationMessage
	device    device.Device
	respondTo *chan int
}

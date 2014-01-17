package notification

import (
	"github.com/gopns/gopns/model"
)

type NotificationTask struct {
	message   NotificationMessage
	device    model.Device
	respondTo *chan int
}

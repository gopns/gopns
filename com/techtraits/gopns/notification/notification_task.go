package notification

import (
	"github.com/gopns/gopns/com/techtraits/gopns/device"
)

type NotificationTask struct {
	message NotificationMessage
	device Device
	respondTo *chan int
}

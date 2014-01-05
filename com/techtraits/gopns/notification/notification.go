package notification

type NotificationMessage struct {
	Title   string
	Message string
}

func (this NotificationMessage) IsValid() bool {
	if len(this.Title) == 0 || len(this.Message) == 0 {
		return false
	} else {
		return true
	}
}
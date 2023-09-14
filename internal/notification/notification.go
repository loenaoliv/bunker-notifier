package notification

type notificationImpl struct {
}

func NewNotification() Notification {
	return &notificationImpl{}
}

type Notification interface {
	Run() error
}

func (j *notificationImpl) Run() error { return nil }

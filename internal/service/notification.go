package service

type NotificationRepository interface {
	SendNotification(string) error
}

package handler

type NotificationChannel interface {
	Send(message string) error
}

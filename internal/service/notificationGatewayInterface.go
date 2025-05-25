package service

type NotificationGateway interface {
	Send(message string) error
}

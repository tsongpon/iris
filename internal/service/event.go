package service

import "time"

type EventRepository interface {
	GetEvents(asOf time.Time) ([]string, error)
	GetEventsBetween(start, end time.Time) ([]string, error)
}

package handler

import "time"

type EventSource interface {
	GetEvents(asOf time.Time) ([]string, error)
}

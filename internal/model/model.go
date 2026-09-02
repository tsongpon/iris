package model

import "time"

type Event struct {
	Name        string
	Description string
	Time        time.Time
}

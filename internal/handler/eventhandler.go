package handler

import (
	"fmt"
	"log"
	"time"
)

type EventHandler struct {
	leaveEventsource EventSource
	notiChannel      NotificationChannel
	asOf             time.Time
}

func NewEventHandler(leaveEventSource EventSource, notiChannel NotificationChannel, asOf time.Time) *EventHandler {
	return &EventHandler{
		leaveEventsource: leaveEventSource,
		notiChannel:      notiChannel,
		asOf:             asOf,
	}
}

func (e *EventHandler) HandleEvent() error {
	leaveEvents, err := e.leaveEventsource.GetEvents(e.asOf)
	if err != nil {
		log.Printf("Failed to get leave events: %v", err)
		return err
	}

	message := fmt.Sprintf("วันนี้ใครลา : (%s)\n", e.asOf.Format(time.DateOnly))
	if len(leaveEvents) == 0 {
		message += "วันนี้ไม่มีคนลา :)"
		log.Printf("No one is on leave today.")
	} else {
		log.Printf("There are " + fmt.Sprint(len(leaveEvents)) + " on leave today.")
		for i, event := range leaveEvents {
			if i == len(leaveEvents)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
	}
	err = e.notiChannel.Send(message)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return err
	}

	return nil
}

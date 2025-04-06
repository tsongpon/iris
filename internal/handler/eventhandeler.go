package handler

import (
	"fmt"
	"log"
	"time"
)

type EventHandler struct {
	leaveEventsource   EventSource
	holidayEventSource EventSource
	notiChannel        NotificationChannel
}

func NewEventHandler(leaveEventSource EventSource, holidayEventSource EventSource, notiChannel NotificationChannel) *EventHandler {
	return &EventHandler{
		leaveEventsource:   leaveEventSource,
		holidayEventSource: holidayEventSource,
		notiChannel:        notiChannel,
	}
}

func (e *EventHandler) HandleEvent() error {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil
	}
	now := time.Now().In(bangkok)

	holidayEvents, err := e.holidayEventSource.GetEvents(now)
	if err != nil {
		log.Printf("Failed to get holiday events: %v", err)
		return err
	}
	leaveEvents, err := e.leaveEventsource.GetEvents(now)
	if err != nil {
		log.Printf("Failed to get leave events: %v", err)
		return err
	}

	if len(holidayEvents) > 0 {
		message := fmt.Sprintf("วันนี้วันหยุด : (%s)\n", now.Format(time.DateOnly))
		log.Printf("Today is a holiday.")
		for i, event := range holidayEvents {
			if i == len(holidayEvents)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
		err = e.notiChannel.Send(message)
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
			return err
		}
	} else {
		message := fmt.Sprintf("วันนี้ใครลา : (%s)\n", now.Format(time.DateOnly))
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
	}
	return nil
}

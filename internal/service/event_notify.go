package service

import (
	"fmt"
	"log"
	"time"
)

type EventNotifyService struct {
	leaveEventRepository   EventRepository
	holidayEventRepository EventRepository
	notificationRepository NotificationRepository
}

func NewEventNotifyService(leaveEventRepo, holidayEventRepo EventRepository, notificationRepo NotificationRepository) EventNotifyService {
	return EventNotifyService{
		leaveEventRepository:   leaveEventRepo,
		holidayEventRepository: holidayEventRepo,
		notificationRepository: notificationRepo,
	}
}

func (e EventNotifyService) Notify(asOf time.Time) error {
	holidayEvents, err := e.holidayEventRepository.GetEvents(asOf)
	if err != nil {
		log.Printf("Error while getting holiday events: %v", err)
		return fmt.Errorf("Error while getting holiday events: %v", err)
	}
	if len(holidayEvents) > 0 {
		log.Println("Today " + asOf.Format(time.DateOnly) + " is a holiday.")
		message := fmt.Sprintf("วันนี้วันหยุด : (%s)\n", asOf.Format(time.DateOnly))
		for i, event := range holidayEvents {
			if i == len(holidayEvents)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
		err = e.notificationRepository.SendNotification(message)
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
			return fmt.Errorf("Error while sending nitification: %v", err)
		}
	} else {
		leaveEvents, err := e.leaveEventRepository.GetEvents(asOf)
		if err != nil {
			return fmt.Errorf("Error while getting events: %v", err)
		}
		if len(leaveEvents) > 0 {
			message := fmt.Sprintf("วันนี้ใครลา : (%s)\n", asOf.Format(time.DateOnly))
			log.Printf("There are " + fmt.Sprint(len(leaveEvents)) + " on leave today.")
			for i, event := range leaveEvents {
				if i == len(leaveEvents)-1 {
					message += fmt.Sprintf("%v", "- "+event)
				} else {
					message += fmt.Sprintf("%v\n", "- "+event)
				}
			}
			err = e.notificationRepository.SendNotification(message)
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
				return fmt.Errorf("Error while sending nitification: %v", err)
			}
		}
	}
	return nil
}

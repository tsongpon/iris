package service

import (
	"fmt"
	"log"
	"time"
)

type LeaveNotifyService struct {
	leaveEventsource EventRepository
	notiGateway      NotificationGateway
	asOf             time.Time
}

func NewLeaveNotifyServicer(leaveEventSource EventRepository, notiGateway NotificationGateway, asOf time.Time) *LeaveNotifyService {
	return &LeaveNotifyService{
		leaveEventsource: leaveEventSource,
		notiGateway:      notiGateway,
		asOf:             asOf,
	}
}

func (e *LeaveNotifyService) HandleEvent() error {
	leaveEvents, err := e.leaveEventsource.GetEvents(e.asOf)
	if err != nil {
		log.Printf("Failed to get leave events: %v", err)
		return err
	}

	if len(leaveEvents) > 0 {
		message := fmt.Sprintf("วันนี้ใครลา : (%s)\n", e.asOf.Format(time.DateOnly))
		log.Printf("There are " + fmt.Sprint(len(leaveEvents)) + " on leave today.")
		for i, event := range leaveEvents {
			if i == len(leaveEvents)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
		err = e.notiGateway.Send(message)
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
			return err
		}
} else {
	log.Printf("No one is on leave today.")
	err = e.notiGateway.Send("No one leave today")
	if err != nil {
		log.Printf("Failed to send 'No one leave today' notification: %v", err)
		return err
	}
	}

	return nil
}

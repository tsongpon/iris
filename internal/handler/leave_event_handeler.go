package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	eventsource "gitlab.com/tsongpon/iris/internal/event-source"
	notichannel "gitlab.com/tsongpon/iris/internal/noti-channel"
)

func LeaveEventHandler() {
	holidayCalendarID := os.Getenv("HOLIDAY_CALENDAR_ID")
	var message string
	holodayEvent, err := eventsource.GetTodayEventFrom(holidayCalendarID)
	if err != nil {
		return
	}
	if len(holodayEvent) > 0 {
		message = fmt.Sprintf("วันนี้วันหยุด : (%s)\n", time.Now().Format(time.DateOnly))
		log.Printf("Today is a holiday.")
		for i, event := range holodayEvent {
			if i == len(holodayEvent)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
	} else {
		leaveCalendatID := os.Getenv("LEAVE_CALENDAR_ID")
		leaveEvents, err := eventsource.GetTodayEventFrom(leaveCalendatID)
		if err != nil {
			return
		}

		message = fmt.Sprintf("วันนี้ใครลา : (%s)\n", time.Now().Format(time.DateOnly))
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
	}

	err = notichannel.SendLineNoti(message)
	if err != nil {
		log.Fatalf("Failed to send LINE notification: %v", err)
	}
}

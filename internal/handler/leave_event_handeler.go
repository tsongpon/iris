package handler

import (
	"fmt"
	"log"
	"time"

	eventsource "gitlab.com/tsongpon/iris/internal/event-source"
	notichannel "gitlab.com/tsongpon/iris/internal/noti-channel"
)

func LeaveEventHandler() {
	events, err := eventsource.GetTodayLeavesEvent()
	if err != nil {
		return
	}

	message := fmt.Sprintf("วันนี้ใครลา : (%s)\n", time.Now().Format(time.DateOnly))
	if len(events) == 0 {
		message += "วันนี้ไม่มีคนลา :)"
		log.Printf("No one is on leave today.")
	} else {
		log.Printf("There are " + fmt.Sprint(len(events)) + " on leave today.")
		for i, event := range events {
			if i == len(events)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
	}

	err = notichannel.SendLineNoti(message)
	if err != nil {
		log.Fatalf("Failed to send LINE notification: %v", err)
	}
}

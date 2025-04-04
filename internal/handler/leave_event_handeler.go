package handler

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
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
		log.Info("No one is on leave today.")
	} else {
		log.Info("There are " + fmt.Sprint(len(events)) + " on leave today.")
		for i, event := range events {
			if i == len(events)-1 {
				message += fmt.Sprintf("%v", "- "+event)
			} else {
				message += fmt.Sprintf("%v\n", "- "+event)
			}
		}
	}

	// Send the message to LINE group
	err = notichannel.SendLineNoti(message)
	if err != nil {
		log.Error("Failed to send LINE notification: ", err)
	}
}

package service

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type EventNotifyService struct {
	leaveEventRepository   EventRepository
	holidayEventRepository EventRepository
	onCallEventRepository  EventRepository
	notificationRepository NotificationRepository
}

func NewEventNotifyService(leaveEventRepo, holidayEventRepo, onCallEventRepo EventRepository,
	notificationRepo NotificationRepository) EventNotifyService {
	return EventNotifyService{
		leaveEventRepository:   leaveEventRepo,
		holidayEventRepository: holidayEventRepo,
		onCallEventRepository:  onCallEventRepo,
		notificationRepository: notificationRepo,
	}
}

func (e EventNotifyService) Notify(asOf time.Time) error {
	if isEndOfMonth(asOf) {
		nextDay := asOf.AddDate(0, 0, 1)
		lastDayOfMonth := time.Date(nextDay.Year(), nextDay.Month()+1, 0, 0, 0, 0, 0, nextDay.Location())
		holidaysNextMonth, err := e.holidayEventRepository.GetEventsBetween(nextDay, lastDayOfMonth)
		if err != nil {
			log.Printf("Error while getting holiday events: %v", err)
			return fmt.Errorf("Error while getting holiday events: %v", err)
		}
		if len(holidaysNextMonth) > 0 {
			log.Println("There are " + strconv.Itoa(len(holidaysNextMonth)) + " holidays next month")
			message := fmt.Sprintf("มีวันหยุด %d วันเดือน %s 🎉🏖️:\n", len(holidaysNextMonth), monthEnToTh(lastDayOfMonth.Format("January")))
			for i, event := range holidaysNextMonth {
				if i == len(holidaysNextMonth)-1 {
					message += fmt.Sprintf("%v", "- "+event)
				} else {
					message += fmt.Sprintf("%v\n", "- "+event)
				}
			}
			err = e.notificationRepository.SendNotification(message)
			if err != nil {
				log.Printf("Error while sending notification: %v", err)
				return fmt.Errorf("Error while sending notification: %v", err)
			}
		} else {
			log.Println("There are no holidays next month")
			message := fmt.Sprintf("เดือน %s ไม่มีวันหยุด 💪😢", monthEnToTh(lastDayOfMonth.Format("January")))
			err = e.notificationRepository.SendNotification(message)
			if err != nil {
				log.Printf("Error while sending notification: %v", err)
				return fmt.Errorf("Error while sending notification: %v", err)
			}
		}
	}

	// Always fetch on-call events (for holidays, weekends, and regular days)
	onCallEvents, err := e.onCallEventRepository.GetEvents(asOf)
	if err != nil {
		log.Printf("Error while getting on-call events: %v", err)
		return fmt.Errorf("Error while getting on-call events: %v", err)
	}

	holidayEvents, err := e.holidayEventRepository.GetEvents(asOf)
	if err != nil {
		log.Printf("Error while getting holiday events: %v", err)
		return fmt.Errorf("Error while getting holiday events: %v", err)
	}

	isWeekend := asOf.Weekday() == time.Saturday || asOf.Weekday() == time.Sunday

	if len(holidayEvents) > 0 || isWeekend {
		message := ""

		if len(holidayEvents) > 0 {
			log.Println("Today " + asOf.Format(time.DateOnly) + " is a holiday.")
			message = fmt.Sprintf("วันนี้วันหยุด 🥳🏖️: (%s)\n", asOf.Format(time.DateOnly))
			for i, event := range holidayEvents {
				if i == len(holidayEvents)-1 {
					message += fmt.Sprintf("%v", "- "+event)
				} else {
					message += fmt.Sprintf("%v\n", "- "+event)
				}
			}
		}

		// Append on-call events on holidays and weekends
		if len(onCallEvents) > 0 {
			log.Printf("There are " + fmt.Sprint(len(onCallEvents)) + " on-call today.")
			if message != "" {
				message += "\n\n"
			}
			message += fmt.Sprintf("📞 วันนี้ใคร On-Call : (%s)\n", asOf.Format(time.DateOnly))
			for i, event := range onCallEvents {
				if i == len(onCallEvents)-1 {
					message += fmt.Sprintf("%v", "- "+event)
				} else {
					message += fmt.Sprintf("%v\n", "- "+event)
				}
			}
		}

		if message != "" {
			err = e.notificationRepository.SendNotification(message)
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
				return fmt.Errorf("Error while sending nitification: %v", err)
			}
		}
	} else {
		leaveEvents, err := e.leaveEventRepository.GetEvents(asOf)
		if err != nil {
			return fmt.Errorf("Error while getting events: %v", err)
		}

		if len(leaveEvents) > 0 || len(onCallEvents) > 0 {
			message := ""

			if len(leaveEvents) > 0 {
				message = fmt.Sprintf("📅 วันนี้ใครลา : (%s)\n", asOf.Format(time.DateOnly))
				log.Printf("There are " + fmt.Sprint(len(leaveEvents)) + " on leave today.")
				for i, event := range leaveEvents {
					if i == len(leaveEvents)-1 {
						message += fmt.Sprintf("%v", "- "+event)
					} else {
						message += fmt.Sprintf("%v\n", "- "+event)
					}
				}
			}

			if len(onCallEvents) > 0 {
				log.Printf("There are " + fmt.Sprint(len(onCallEvents)) + " on-call today.")
				if message != "" {
					message += "\n\n"
				}
				message += fmt.Sprintf("📞 วันนี้ใคร On-Call : (%s)\n", asOf.Format(time.DateOnly))
				for i, event := range onCallEvents {
					if i == len(onCallEvents)-1 {
						message += fmt.Sprintf("%v", "- "+event)
					} else {
						message += fmt.Sprintf("%v\n", "- "+event)
					}
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

func isEndOfMonth(date time.Time) bool {
	// Add one day to the date and check if the month changes
	nextDay := date.AddDate(0, 0, 1)
	return nextDay.Month() != date.Month()
}

func monthEnToTh(monthEn string) string {
	switch monthEn {
	case "January":
		return "มกราคม"
	case "February":
		return "กุมภาพันธ์"
	case "March":
		return "มีนาคม"
	case "April":
		return "เมษายน"
	case "May":
		return "พฤษภาคม"
	case "June":
		return "มิถุนายน"
	case "July":
		return "กรกฎาคม"
	case "August":
		return "สิงหาคม"
	case "September":
		return "กันยายน"
	case "October":
		return "ตุลาคม"
	case "November":
		return "พฤศจิกายน"
	case "December":
		return "ธันวาคม"
	default:
		return monthEn
	}
}

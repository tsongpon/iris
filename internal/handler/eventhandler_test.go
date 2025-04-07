package handler

import (
	"fmt"
	"testing"
	"time"
)

type MockEventSource struct {
	event []string
	err   error
}

func (m *MockEventSource) GetEvents(asOf time.Time) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.event, nil
}

type MockNotificationChannel struct {
	numberOfCalls int
	sentMessage   string
}

func (m *MockNotificationChannel) Send(message string) error {
	m.numberOfCalls++
	m.sentMessage = message
	return nil
}

func TestEventHandlerWith2LeaveNoHoliday(t *testing.T) {
	mockNotificationChannel := &MockNotificationChannel{}
	mockLeaveEventSource := &MockEventSource{event: []string{"Tum leave", "Songpon leave"}}
	mockHolidayEventSource := &MockEventSource{event: []string{}}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)

	handler := NewEventHandler(mockLeaveEventSource, mockHolidayEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if mockNotificationChannel.numberOfCalls != 1 {
		t.Errorf("Expected notification channel to be called once, but got %d", mockNotificationChannel.numberOfCalls)
	}

	if mockNotificationChannel.sentMessage == "" {
		t.Errorf("Expected a message to be sent, but got an empty message")
	}

	if mockNotificationChannel.sentMessage != "วันนี้ใครลา : ("+now.Format(time.DateOnly)+")\n- Tum leave\n- Songpon leave" {
		t.Errorf("Expected message to be 'วันนี้ใครลา : (%s)\n- Tum leave\n- Songpon leave', but got %s", now.Format(time.DateOnly), mockNotificationChannel.sentMessage)
	}
}

func TestEventHandlerWith1LeaveAnd1Holiday(t *testing.T) {
	mockNotificationChannel := &MockNotificationChannel{}
	mockLeaveEventSource := &MockEventSource{event: []string{"Tum leave"}}
	mockHolidayEventSource := &MockEventSource{event: []string{"Father's Day"}}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewEventHandler(mockLeaveEventSource, mockHolidayEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if mockNotificationChannel.numberOfCalls != 1 {
		t.Errorf("Expected notification channel to be called once, but got %d", mockNotificationChannel.numberOfCalls)
	}

	if mockNotificationChannel.sentMessage == "" {
		t.Errorf("Expected a message to be sent, but got an empty message")
	}

	if mockNotificationChannel.sentMessage != "วันนี้วันหยุด : ("+now.Format(time.DateOnly)+")\n- Father's Day" {
		t.Errorf("Expected message to be 'วันนี้วันหยุด : (%s)\n- Father's Day', but got %s", now.Format(time.DateOnly), mockNotificationChannel.sentMessage)
	}
}

func TestEventHandlerWithNoLeaveAndNoHoliday(t *testing.T) {
	mockNotificationChannel := &MockNotificationChannel{}
	mockLeaveEventSource := &MockEventSource{event: []string{}}
	mockHolidayEventSource := &MockEventSource{event: []string{}}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewEventHandler(mockLeaveEventSource, mockHolidayEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if mockNotificationChannel.numberOfCalls != 1 {
		t.Errorf("Expected notification channel to be called once, but got %d", mockNotificationChannel.numberOfCalls)
	}

	if mockNotificationChannel.sentMessage == "" {
		t.Errorf("Expected a message to be sent, but got an empty message")
	}

	if mockNotificationChannel.sentMessage != "วันนี้ใครลา : ("+now.Format(time.DateOnly)+")\nวันนี้ไม่มีคนลา :)" {
		t.Errorf("Expected message to be 'วันนี้ใครลา : (%s)\nวันนี้ไม่มีคนลา :)', but got %s", now.Format(time.DateOnly), mockNotificationChannel.sentMessage)
	}
}

func TestEventHandlerWithHolidayEventSourceError(t *testing.T) {
	mockNotificationChannel := &MockNotificationChannel{}
	mockLeaveEventSource := &MockEventSource{event: []string{}}
	mockHolidayEventSource := &MockEventSource{event: []string{}, err: fmt.Errorf("holiday event source error")}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewEventHandler(mockLeaveEventSource, mockHolidayEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	if mockNotificationChannel.numberOfCalls != 0 {
		t.Errorf("Expected notification channel to not be called, but got %d", mockNotificationChannel.numberOfCalls)
	}
}

func TestEventHandlerWithLeaveEventSourceError(t *testing.T) {
	mockNotificationChannel := &MockNotificationChannel{}
	mockHolidayEventSource := &MockEventSource{event: []string{}}
	mockLeaveEventSource := &MockEventSource{event: []string{}, err: fmt.Errorf("leave event source error")}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewEventHandler(mockLeaveEventSource, mockHolidayEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	if mockNotificationChannel.numberOfCalls != 0 {
		t.Errorf("Expected notification channel to not be called, but got %d", mockNotificationChannel.numberOfCalls)
	}
}

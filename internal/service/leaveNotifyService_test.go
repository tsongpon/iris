package service

import (
	"fmt"
	"testing"
	"time"
)

type MockEventRepository struct {
	event []string
	err   error
}

func (m *MockEventRepository) GetEvents(asOf time.Time) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.event, nil
}

type MockNotificationGateway struct {
	numberOfCalls int
	sentMessage   string
}

func (m *MockNotificationGateway) Send(message string) error {
	m.numberOfCalls++
	m.sentMessage = message
	return nil
}

func TestEventHandlerWith2Leaves(t *testing.T) {
	mockNotificationChannel := &MockNotificationGateway{}
	mockLeaveEventSource := &MockEventRepository{event: []string{"Tum leave", "Songpon leave"}}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)

	handler := NewLeaveNotifyServicer(mockLeaveEventSource, mockNotificationChannel, now)
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

func TestEventHandlerWithNoLeave(t *testing.T) {
	mockNotificationChannel := &MockNotificationGateway{}
	mockLeaveEventSource := &MockEventRepository{event: []string{}}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewLeaveNotifyServicer(mockLeaveEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if mockNotificationChannel.numberOfCalls != 0 {
		t.Errorf("Expected notification channel not to be called, but got %d", mockNotificationChannel.numberOfCalls)
	}
}

func TestEventHandlerWithLeaveEventSourceError(t *testing.T) {
	mockNotificationChannel := &MockNotificationGateway{}
	mockLeaveEventSource := &MockEventRepository{event: []string{}, err: fmt.Errorf("leave event source error")}

	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(bangkok)
	handler := NewLeaveNotifyServicer(mockLeaveEventSource, mockNotificationChannel, now)
	err := handler.HandleEvent()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	if mockNotificationChannel.numberOfCalls != 0 {
		t.Errorf("Expected notification channel to not be called, but got %d", mockNotificationChannel.numberOfCalls)
	}
}

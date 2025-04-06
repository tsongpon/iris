package handler

import (
	"testing"
	"time"
)

type MockLeaveEventSource struct{}

func (m *MockLeaveEventSource) GetEvents(asOf time.Time) ([]string, error) {
	return []string{"Mock Event 1", "Mock Event 2"}, nil
}

type MockHolidayEventSource struct{}

func (m *MockHolidayEventSource) GetEvents(asOf time.Time) ([]string, error) {
	return []string{}, nil
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

func TestEventHandler(t *testing.T) {
	mockNotificationChennel := &MockNotificationChannel{}
	handler := NewEventHandler(&MockLeaveEventSource{}, &MockHolidayEventSource{}, mockNotificationChennel)
	err := handler.HandleEvent()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if mockNotificationChennel.numberOfCalls != 1 {
		t.Errorf("Expected notification channel to be called once, but got %d", mockNotificationChennel.numberOfCalls)
	}

	if mockNotificationChennel.sentMessage == "" {
		t.Errorf("Expected a message to be sent, but got an empty message")
	}
	if mockNotificationChennel.sentMessage != "วันนี้ใครลา : ("+time.Now().In(time.FixedZone("Asia/Bangkok", 7*3600)).Format(time.DateOnly)+")\n- Mock Event 1\n- Mock Event 2" {
		t.Errorf("Expected message to be 'วันนี้ใครลา : (%s)\n- Mock Event 1\n- Mock Event 2', but got %s", time.Now().In(time.FixedZone("Asia/Bangkok", 7*3600)).Format(time.DateOnly), mockNotificationChennel.sentMessage)
	}
}

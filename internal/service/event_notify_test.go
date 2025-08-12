package service

import (
	"errors"
	"testing"
	"time"
)

// Mock implementations
type MockEventRepository struct {
	events []string
	err    error
}

func (m *MockEventRepository) GetEvents(asOf time.Time) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

type MockNotificationRepository struct {
	numberOfCalls int
	sentMessage   string
	err           error
}

func (m *MockNotificationRepository) SendNotification(message string) error {
	m.numberOfCalls++
	m.sentMessage = message
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestEventNotifyService_Notify_HolidayEvents_SingleEvent(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{"National Day"}}
	mockLeaveRepo := &MockEventRepository{events: []string{}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "วันนี้วันหยุด 🎉🏖️: (2025-08-12)\n- National Day"
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_HolidayEvents_MultipleEvents(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{"National Day", "Independence Day"}}
	mockLeaveRepo := &MockEventRepository{events: []string{}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "วันนี้วันหยุด 🎉🏖️: (2025-08-12)\n- National Day\n- Independence Day"
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_HolidayEvents_GetEventsError(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{err: errors.New("holiday repository error")}
	mockLeaveRepo := &MockEventRepository{events: []string{}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "Error while getting holiday events: holiday repository error"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if mockNotification.numberOfCalls != 0 {
		t.Errorf("Expected notification not to be called, got %d calls", mockNotification.numberOfCalls)
	}
}

func TestEventNotifyService_Notify_HolidayEvents_SendNotificationError(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{err: errors.New("notification error")}
	mockHolidayRepo := &MockEventRepository{events: []string{"National Day"}}
	mockLeaveRepo := &MockEventRepository{events: []string{}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "Error while sending nitification: notification error"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}
}

func TestEventNotifyService_Notify_LeaveEvents_SingleEvent(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{events: []string{"John Doe"}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "📅 วันนี้ใครลา : (2025-08-12)\n- John Doe"
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_LeaveEvents_MultipleEvents(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{events: []string{"John Doe", "Jane Smith", "Bob Johnson"}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "📅 วันนี้ใครลา : (2025-08-12)\n- John Doe\n- Jane Smith\n- Bob Johnson"
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_LeaveEvents_GetEventsError(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{err: errors.New("leave repository error")}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "Error while getting events: leave repository error"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if mockNotification.numberOfCalls != 0 {
		t.Errorf("Expected notification not to be called, got %d calls", mockNotification.numberOfCalls)
	}
}

func TestEventNotifyService_Notify_LeaveEvents_SendNotificationError(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{err: errors.New("notification error")}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{events: []string{"John Doe"}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "Error while sending nitification: notification error"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}
}

func TestEventNotifyService_Notify_NoEventsAtAll(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{events: []string{}}   // No leaves

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 0 {
		t.Errorf("Expected notification not to be called, got %d calls", mockNotification.numberOfCalls)
	}

	if mockNotification.sentMessage != "" {
		t.Errorf("Expected no message to be sent, got '%s'", mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_HolidayTakesPrecedenceOverLeave(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{"National Day"}}
	mockLeaveRepo := &MockEventRepository{events: []string{"John Doe"}} // This should be ignored

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	// Should send holiday message, not leave message
	expectedMessage := "วันนี้วันหยุด 🎉🏖️: (2025-08-12)\n- National Day"
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected holiday message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_EmptyHolidayStringInSlice(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{""}} // Empty string but slice is not empty
	mockLeaveRepo := &MockEventRepository{events: []string{}}

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "วันนี้วันหยุด 🎉🏖️: (2025-08-12)\n- "
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

func TestEventNotifyService_Notify_EmptyLeaveStringInSlice(t *testing.T) {
	// Arrange
	mockNotification := &MockNotificationRepository{}
	mockHolidayRepo := &MockEventRepository{events: []string{}} // No holidays
	mockLeaveRepo := &MockEventRepository{events: []string{""}} // Empty string but slice is not empty

	service := NewEventNotifyService(mockLeaveRepo, mockHolidayRepo, mockNotification)

	// Create the specific date: 12 August 2025, 8 AM Bangkok time
	bangkok, _ := time.LoadLocation("Asia/Bangkok")
	testDate := time.Date(2025, 8, 12, 8, 0, 0, 0, bangkok)

	// Act
	err := service.Notify(testDate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockNotification.numberOfCalls != 1 {
		t.Errorf("Expected notification to be called once, got %d", mockNotification.numberOfCalls)
	}

	expectedMessage := "📅 วันนี้ใครลา : (2025-08-12)\n- "
	if mockNotification.sentMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, mockNotification.sentMessage)
	}
}

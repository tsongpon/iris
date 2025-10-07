package service

import (
	"errors"
	"testing"
	"time"
)

// Mock implementations
type MockEventRepository struct {
	events        []string
	eventsBetween []string
	err           error
}

func (m *MockEventRepository) GetEvents(asOf time.Time) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func (m *MockEventRepository) GetEventsBetween(start, end time.Time) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.eventsBetween, nil
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

// Tests for isEndOfMonth function
func TestIsEndOfMonth_LastDayOfJanuary(t *testing.T) {
	// Test January 31st (31-day month)
	date := time.Date(2024, 1, 31, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for January 31st, got false")
	}
}

func TestIsEndOfMonth_NotLastDayOfJanuary(t *testing.T) {
	// Test January 30th (not last day)
	date := time.Date(2024, 1, 30, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if result {
		t.Errorf("Expected false for January 30th, got true")
	}
}

func TestIsEndOfMonth_LastDayOfFebruaryLeapYear(t *testing.T) {
	// Test February 29th in a leap year
	date := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for February 29th (leap year), got false")
	}
}

func TestIsEndOfMonth_LastDayOfFebruaryNonLeapYear(t *testing.T) {
	// Test February 28th in a non-leap year
	date := time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for February 28th (non-leap year), got false")
	}
}

func TestIsEndOfMonth_NotLastDayOfFebruaryLeapYear(t *testing.T) {
	// Test February 28th in a leap year (not last day)
	date := time.Date(2024, 2, 28, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if result {
		t.Errorf("Expected false for February 28th (leap year), got true")
	}
}

func TestIsEndOfMonth_LastDayOfApril(t *testing.T) {
	// Test April 30th (30-day month)
	date := time.Date(2024, 4, 30, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for April 30th, got false")
	}
}

func TestIsEndOfMonth_NotLastDayOfApril(t *testing.T) {
	// Test April 29th (not last day)
	date := time.Date(2024, 4, 29, 12, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if result {
		t.Errorf("Expected false for April 29th, got true")
	}
}

func TestIsEndOfMonth_LastDayOfDecember(t *testing.T) {
	// Test December 31st (year end)
	date := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for December 31st, got false")
	}
}

func TestIsEndOfMonth_FirstDayOfMonth(t *testing.T) {
	// Test first day of month
	date := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if result {
		t.Errorf("Expected false for March 1st, got true")
	}
}

func TestIsEndOfMonth_MiddleOfMonth(t *testing.T) {
	// Test middle of month
	date := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	result := isEndOfMonth(date)
	if result {
		t.Errorf("Expected false for June 15th, got true")
	}
}

func TestIsEndOfMonth_WithBangkokTimezone(t *testing.T) {
	// Test with Bangkok timezone
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatalf("Failed to load Bangkok timezone: %v", err)
	}
	date := time.Date(2024, 3, 31, 20, 0, 0, 0, bangkok)
	result := isEndOfMonth(date)
	if !result {
		t.Errorf("Expected true for March 31st in Bangkok timezone, got false")
	}
}

func TestIsEndOfMonth_AllMonths(t *testing.T) {
	// Test all 12 months with their respective last days
	testCases := []struct {
		month   time.Month
		lastDay int
	}{
		{time.January, 31},
		{time.February, 28}, // Non-leap year
		{time.March, 31},
		{time.April, 30},
		{time.May, 31},
		{time.June, 30},
		{time.July, 31},
		{time.August, 31},
		{time.September, 30},
		{time.October, 31},
		{time.November, 30},
		{time.December, 31},
	}

	for _, tc := range testCases {
		date := time.Date(2023, tc.month, tc.lastDay, 12, 0, 0, 0, time.UTC)
		result := isEndOfMonth(date)
		if !result {
			t.Errorf("Expected true for last day of %s (day %d), got false", tc.month, tc.lastDay)
		}

		// Also test the day before (should be false)
		dateBefore := time.Date(2023, tc.month, tc.lastDay-1, 12, 0, 0, 0, time.UTC)
		resultBefore := isEndOfMonth(dateBefore)
		if resultBefore {
			t.Errorf("Expected false for %s %d (day before last), got true", tc.month, tc.lastDay-1)
		}
	}
}

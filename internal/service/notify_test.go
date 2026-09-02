package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/tsongpon/iris/internal/model"
)

type fakeEventRepo struct {
	events []model.Event
	err    error
}

func (r *fakeEventRepo) GetEvents(ctx context.Context, asOf time.Time) ([]model.Event, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.events, nil
}

type fakeNotificationRepo struct {
	called   int
	messages []string
	err      error
}

func (r *fakeNotificationRepo) SendNotification(ctx context.Context, message string) error {
	r.called++
	r.messages = append(r.messages, message)
	return r.err
}

func TestSendEventNotificationSendsEventList(t *testing.T) {
	asOf := time.Date(2026, 9, 2, 15, 0, 0, 0, time.UTC)
	events := []model.Event{
		{Name: "Alice (Annual Leave)", Time: asOf},
		{Name: "Bob (Sick Leave)", Time: asOf},
		{Name: "Team Standup", Time: asOf},
	}
	eventRepo := &fakeEventRepo{events: events}
	notificationRepo := &fakeNotificationRepo{}

	svc, err := NewNotifyService(notificationRepo, eventRepo)
	if err != nil {
		t.Fatalf("NewNotifyService: %v", err)
	}
	if err := svc.SendNotify(context.Background(), asOf); err != nil {
		t.Fatalf("SendEventNotification: %v", err)
	}
	if notificationRepo.called != 1 {
		t.Errorf("SendNotification called %d times, want 1", notificationRepo.called)
	}
	msg := notificationRepo.messages[0]
	for _, name := range []string{"Alice (Annual Leave)", "Bob (Sick Leave)", "Team Standup"} {
		if !strings.Contains(msg, name) {
			t.Errorf("message missing event %q:\n%s", name, msg)
		}
	}
	if !strings.Contains(msg, "📅 Events on Wednesday, 2 September 2026") {
		t.Errorf("message missing dated title:\n%s", msg)
	}
}

func TestSendEventNotificationSkipsEmptyDay(t *testing.T) {
	eventRepo := &fakeEventRepo{events: nil}
	notificationRepo := &fakeNotificationRepo{}

	svc, err := NewNotifyService(notificationRepo, eventRepo)
	if err != nil {
		t.Fatalf("NewNotifyService: %v", err)
	}
	if err := svc.SendNotify(context.Background(), time.Now()); err != nil {
		t.Fatalf("SendEventNotification: %v", err)
	}
	if notificationRepo.called != 0 {
		t.Errorf("SendNotification called %d times, want 0", notificationRepo.called)
	}
}

func TestSendEventNotificationGetEventsError(t *testing.T) {
	eventRepo := &fakeEventRepo{err: errors.New("calendar down")}
	notificationRepo := &fakeNotificationRepo{}

	svc, err := NewNotifyService(notificationRepo, eventRepo)
	if err != nil {
		t.Fatalf("NewNotifyService: %v", err)
	}
	err = svc.SendNotify(context.Background(), time.Now())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "calendar down") {
		t.Errorf("error should wrap cause: %v", err)
	}
	if notificationRepo.called != 0 {
		t.Errorf("SendNotification called %d times, want 0", notificationRepo.called)
	}
}

func TestSendEventNotificationSendError(t *testing.T) {
	events := []model.Event{{Name: "Alice (Annual Leave)", Time: time.Now()}}
	eventRepo := &fakeEventRepo{events: events}
	notificationRepo := &fakeNotificationRepo{err: errors.New("webhook failed")}

	svc, err := NewNotifyService(notificationRepo, eventRepo)
	if err != nil {
		t.Fatalf("NewNotifyService: %v", err)
	}
	err = svc.SendNotify(context.Background(), time.Now())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "webhook failed") {
		t.Errorf("error should wrap cause: %v", err)
	}
}

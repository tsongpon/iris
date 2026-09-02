package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/tsongpon/iris/internal/model"
)

type EventRepository interface {
	GetEvents(ctx context.Context, asOf time.Time) ([]model.Event, error)
}

type NotificationRepository interface {
	SendNotification(ctx context.Context, message string) error
}

type NotifyService struct {
	notificationRepo NotificationRepository
	eventRepo        EventRepository
}

func NewNotifyService(notificationRepo NotificationRepository, eventRepo EventRepository) (*NotifyService, error) {
	return &NotifyService{
		notificationRepo: notificationRepo,
		eventRepo:        eventRepo,
	}, nil
}

func (s *NotifyService) SendNotify(ctx context.Context, asOf time.Time) error {
	events, err := s.eventRepo.GetEvents(ctx, asOf)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get events for notification",
			"asOf", asOf.Format("2006-01-02"),
			"err", err,
		)
		return fmt.Errorf("get events: %w", err)
	}
	if len(events) == 0 {
		return nil
	}

	lines := make([]string, 0, len(events))
	for _, e := range events {
		lines = append(lines, "• "+e.Name)
	}
	message := "📅 Events on " + asOf.Format("Monday, 2 January 2006") + "\n\n" +
		strings.Join(lines, "\n")

	if err := s.notificationRepo.SendNotification(ctx, message); err != nil {
		slog.ErrorContext(ctx, "failed to send event notification",
			"asOf", asOf.Format("2006-01-02"),
			"eventCount", len(events),
			"err", err,
		)
		return fmt.Errorf("send notification: %w", err)
	}

	slog.InfoContext(ctx, "event notification sent",
		"asOf", asOf.Format("2006-01-02"),
		"eventCount", len(events),
	)
	return nil
}

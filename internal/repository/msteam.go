package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tsongpon/iris/internal/service"
)

var _ service.NotificationRepository = (*MSTeamNotificationRepository)(nil)

type MSTeamNotificationRepository struct {
	webhookURL string
}

func NewMSTeamNotificationRepository(webhookURL string) *MSTeamNotificationRepository {
	return &MSTeamNotificationRepository{webhookURL: webhookURL}
}

func (r *MSTeamNotificationRepository) SendNotification(ctx context.Context, message string) error {
	payload := map[string]any{
		"type": "message",
		"attachments": []map[string]any{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]any{
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type":    "AdaptiveCard",
					"version": "1.4",
					"body": []map[string]any{
						{
							"type":     "TextBlock",
							"text":     message,
							"wrap":     true,
							"fontType": "Default",
							"size":     "Medium",
							"spacing":  "Small",
						},
					},
				},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Make the POST request
	webhookURL := r.webhookURL

	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send notification", "error", err)
		return err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send notification", "error", err)
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		slog.ErrorContext(ctx, "failed to send notification", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("failed to send notification: %s", string(body))
	}

	return nil
}

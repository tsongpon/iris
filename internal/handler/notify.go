package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/tsongpon/iris/internal/service"
)

type NotifyHandler struct {
	notifyService *service.NotifyService
}

func NewNotifyHandler(notifyService *service.NotifyService) (*NotifyHandler, error) {
	return &NotifyHandler{
		notifyService: notifyService,
	}, nil
}

func (h *NotifyHandler) SendNotify(c *echo.Context) error {
	asOf := time.Now()
	slog.Info("Sending notify as of", "asOf", asOf)
	err := h.notifyService.SendNotify(c.Request().Context(), asOf)
	if err != nil {
		slog.Error("Failed to send notify", "error", err)
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

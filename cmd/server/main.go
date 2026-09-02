package main

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"

	"github.com/tsongpon/iris/internal/handler"
	"github.com/tsongpon/iris/internal/repository"
	"github.com/tsongpon/iris/internal/service"
)

func main() {
	// Load .env if present. Values already set in the real environment are
	// never overridden, so container/production config always wins.
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			slog.Warn("load .env failed, continuing with environment only", "err", err)
		}
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(),
	}))
	slog.SetDefault(logger)

	// Credentials are optional: if neither GOOGLE_CREDENTIALS_JSON nor
	// GOOGLE_CREDENTIALS_FILE is set, Application Default Credentials are
	// used — the natural fit for containers on GCP (Workload Identity,
	// Cloud Run, GKE, GCE) where the metadata server provides tokens.
	var credentialsJSON []byte
	var err error
	if raw := os.Getenv("GOOGLE_CREDENTIALS_JSON"); raw != "" {
		credentialsJSON = []byte(raw)
	} else if path := os.Getenv("GOOGLE_CREDENTIALS_FILE"); path != "" {
		credentialsJSON, err = os.ReadFile(path)
		if err != nil {
			slog.Error("read credentials file failed", "path", path, "err", err)
			os.Exit(1)
		}
	}

	calendarID := os.Getenv("GOOGLE_CALENDAR_ID")
	if calendarID == "" {
		slog.Error("GOOGLE_CALENDAR_ID env var is required")
		os.Exit(1)
	}

	timezone := os.Getenv("CALENDAR_TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Bangkok"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error("load location failed", "timezone", timezone, "err", err)
		os.Exit(1)
	}

	eventRepo, err := repository.NewGoogleCalendarEventRepository(
		context.Background(),
		credentialsJSON,
		calendarID,
		loc,
	)
	if err != nil {
		slog.Error("create repository failed", "err", err)
		os.Exit(1)
	}

	msTeamWebhookURL := os.Getenv("MSTEAM_WEBHOOK_URL")
	if msTeamWebhookURL == "" {
		slog.Error("MSTEAM_WEBHOOK_URL env var is required")
		os.Exit(1)
	}
	notificationRepo := repository.NewMSTeamNotificationRepository(msTeamWebhookURL)

	notifyService, err := service.NewNotifyService(notificationRepo, eventRepo)
	if err != nil {
		slog.Error("create notify service failed", "err", err)
		os.Exit(1)
	}
	notifyHandler, err := handler.NewNotifyHandler(notifyService)
	if err != nil {
		slog.Error("create notify handler failed", "err", err)
		os.Exit(1)
	}

	e := echo.New()
	e.GET("/ping", func(c *echo.Context) error {
		return c.String(200, "pong")
	})
	e.POST("v1/notify", notifyHandler.SendNotify)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	if err := e.Start(":" + port); err != nil {
		slog.Error("server start failed", "err", err)
		os.Exit(1)
	}
}

// logLevel maps the LOG_LEVEL env var to a slog level, defaulting to info.
func logLevel() slog.Level {
	switch s := os.Getenv("LOG_LEVEL"); s {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

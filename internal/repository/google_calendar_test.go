package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)


func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func TestGetEventsReturnsEventsInWindow(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	asOf := time.Date(2026, 9, 1, 15, 30, 0, 0, loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/test@example.com/events" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("timeMin") != "2026-09-01T00:00:00+07:00" {
			t.Errorf("timeMin = %q", q.Get("timeMin"))
		}
		if q.Get("timeMax") != "2026-09-02T00:00:00+07:00" {
			t.Errorf("timeMax = %q", q.Get("timeMax"))
		}
		if q.Get("singleEvents") != "true" {
			t.Errorf("singleEvents = %q", q.Get("singleEvents"))
		}
		writeJSON(t, w, &calendar.Events{
			Items: []*calendar.Event{
				{
					Summary:     "Standup",
					Description: "daily sync",
					Status:      "confirmed",
					Start:       &calendar.EventDateTime{DateTime: "2026-09-01T09:00:00+07:00"},
				},
				{
					Summary: "Cancelled event",
					Status:  "cancelled",
					Start:   &calendar.EventDateTime{DateTime: "2026-09-01T12:00:00+07:00"},
				},
			},
		})
	}))
	defer ts.Close()

	srv, err := calendar.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("create calendar service: %v", err)
	}
	repo := &GoogleCalendarEventRepository{srv: srv, calendarID: "test@example.com", loc: loc}

	events, err := repo.GetEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1: %+v", len(events), events)
	}
	if events[0].Name != "Standup" {
		t.Errorf("Name = %q", events[0].Name)
	}
	if events[0].Description != "daily sync" {
		t.Errorf("Description = %q", events[0].Description)
	}
	want, _ := time.Parse(time.RFC3339, "2026-09-01T09:00:00+07:00")
	if !events[0].Time.Equal(want) {
		t.Errorf("Time = %v, want %v", events[0].Time, want)
	}
}

func TestGetEventsPaginates(t *testing.T) {
	loc := time.UTC
	asOf := time.Date(2026, 9, 1, 8, 0, 0, 0, loc)
	var calls int

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if q := r.URL.Query().Get("pageToken"); q == "" {
			writeJSON(t, w, &calendar.Events{
				NextPageToken: "page2",
				Items: []*calendar.Event{{
					Summary: "Morning",
					Start:   &calendar.EventDateTime{DateTime: "2026-09-01T09:00:00Z"},
				}},
			})
			return
		} else if q != "page2" {
			t.Errorf("pageToken = %q, want page2", q)
		}
		writeJSON(t, w, &calendar.Events{
			Items: []*calendar.Event{{
				Summary: "Afternoon",
				Start:   &calendar.EventDateTime{DateTime: "2026-09-01T14:00:00Z"},
			}},
		})
	}))
	defer ts.Close()

	srv, err := calendar.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("create calendar service: %v", err)
	}
	repo := &GoogleCalendarEventRepository{srv: srv, calendarID: "test@example.com", loc: loc}

	events, err := repo.GetEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if calls != 2 {
		t.Errorf("calls = %d, want 2", calls)
	}
	if len(events) != 2 {
		t.Fatalf("got %d events, want 2: %+v", len(events), events)
	}
	if events[0].Name != "Morning" || events[1].Name != "Afternoon" {
		t.Errorf("order = %q, %q", events[0].Name, events[1].Name)
	}
}

func TestGetEventsAllDayEvent(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	asOf := time.Date(2026, 9, 1, 23, 0, 0, 0, loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, &calendar.Events{
			Items: []*calendar.Event{{
				Summary: "Conference",
				Start:   &calendar.EventDateTime{Date: "2026-09-01"},
			}},
		})
	}))
	defer ts.Close()

	srv, err := calendar.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("create calendar service: %v", err)
	}
	repo := &GoogleCalendarEventRepository{srv: srv, calendarID: "test@example.com", loc: loc}

	events, err := repo.GetEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1: %+v", len(events), events)
	}
	want := time.Date(2026, 9, 1, 0, 0, 0, 0, loc)
	if !events[0].Time.Equal(want) {
		t.Errorf("Time = %v, want %v", events[0].Time, want)
	}
}

func TestNewGoogleCalendarEventRepository(t *testing.T) {
	ctx := context.Background()
	loc := time.UTC

	if _, err := NewGoogleCalendarEventRepository(ctx, nil, "", loc); err == nil {
		t.Error("expected error for empty calendar id")
	}
}

// With empty credentials the constructor must fall back to Application
// Default Credentials. On machines without ADC configured, construction may
// fail with the ADC "could not find default credentials" error — that is
// acceptable; any other error (e.g. rejecting empty credentials) is not.
func TestNewGoogleCalendarEventRepositoryADCFallback(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("HOME", t.TempDir())

	_, err := NewGoogleCalendarEventRepository(context.Background(), nil, "test@example.com", time.UTC)
	if err == nil {
		return
	}
	if !strings.Contains(err.Error(), "could not find default credentials") {
		t.Errorf("expected ADC fallback (or its missing-credentials error), got: %v", err)
	}
}

func TestGetEventsAPIError(t *testing.T) {
	loc := time.UTC
	asOf := time.Date(2026, 9, 1, 0, 0, 0, 0, loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		writeJSON(t, w, map[string]any{
			"error": map[string]any{"code": 403, "message": "forbidden"},
		})
	}))
	defer ts.Close()

	srv, err := calendar.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("create calendar service: %v", err)
	}
	repo := &GoogleCalendarEventRepository{srv: srv, calendarID: "test@example.com", loc: loc}

	if _, err := repo.GetEvents(context.Background(), asOf); err == nil {
		t.Fatal("GetEvents: expected error, got nil")
	}
}
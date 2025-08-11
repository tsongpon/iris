# Iris

![Deploy Lambda](https://github.com/tsongpon/iris/actions/workflows/deploy.yml/badge.svg)

A Go application that reads events from Google Calendar and sends notifications to Line groups. The system monitors both holiday events and leave events, providing automated notifications for team coordination.

## Features

- **Google Calendar Integration**: Fetches events from multiple Google Calendars
- **Line Messaging**: Sends automated notifications to Line groups
- **Holiday Detection**: Prioritizes holiday notifications over leave notifications

## Architecture

The project follows clean architecture principles with the Repository Pattern:

```
iris/
├── cmd/
│   └── iris/           # Application entry point
├── internal/
│   ├── repository/     # Data access layer
│   │   ├── google_calendar.go
│   │   └── line_notification.go
│   └── service/        # Business logic layer
│       ├── event.go
│       ├── event_notify.go
│       ├── event_notify_test.go
│       └── notification.go
├── pkg/                # Shared packages
├── Dockerfile          # Container configuration
├── go.mod              # Go module dependencies
└── README.md
```

### Key Components

- **EventNotifyService**: Core business logic for event notification
- **GoogleCalendar**: Repository for Google Calendar API integration
- **LineNotificationRepository**: Repository for Line messaging API
- **EventRepository Interface**: Abstraction for event data sources
- **NotificationRepository Interface**: Abstraction for notification channels

## Prerequisites

- Go 1.23.4 or later
- Google Calendar API credentials (Service Account)
- Line Bot API credentials
- Access to Google Calendars you want to monitor

## Installation

1. Clone the repository:
```bash
git clone https://github.com/tsongpon/iris.git
cd iris
```

2. Install dependencies:
```bash
go mod tidy
```

## Configuration

Set the following environment variables:

```bash
# Google Calendar Configuration
export GOOGLE_CREDENTIALS_JSON=$(cat your-google-calendar-credential-file | base64)
export LEAVE_CALENDAR_ID=your-leave-calendar-id
export HOLIDAY_CALENDAR_ID=en.th#holiday@group.v.calendar.google.com

# Line Bot Configuration
export LINE_CHANNEL_SECRET=your-line-channel-secret
export LINE_CHANNEL_TOKEN=your-line-channel-token
export LINE_GROUP_ID=your-line-group-id-to-send-message-to
```

### Google Calendar Setup

1. Create a Google Cloud Project
2. Enable the Google Calendar API
3. Create a Service Account and download the JSON credentials
4. Share your Google Calendar with the Service Account email
5. Base64 encode the credentials JSON file

### Line Bot Setup

1. Create a Line Developer account
2. Create a new Line Bot
3. Get the Channel Secret and Channel Token
4. Add the bot to your Line group
5. Get the Group ID where notifications should be sent

## Usage

### Running Locally

```bash
go run cmd/iris/main.go
```

### Running with Docker

```bash
docker build -t iris .
docker run --env-file .env iris
```

### Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run service layer tests specifically:
```bash
go test ./internal/service -v
```

The test suite includes:
- **12 comprehensive unit tests** covering all scenarios
- **Error handling tests** for all failure cases
- **Business logic validation** including message formatting

## How It Works

1. **Event Fetching**: The application fetches events from two Google Calendars:
   - Holiday calendar (Thai holidays)
   - Leave calendar (employee leave requests)

2. **Event Processing**: For a specific date (currently set to August 12, 2025, 8 AM Bangkok time):
   - First checks for holiday events
   - If holidays exist, sends holiday notification and ignores leave events
   - If no holidays, checks for leave events and sends leave notification
   - If no events at all, sends no notification

3. **Notification Format**:
   - **Holiday**: `วันนี้วันหยุด : (2025-08-12)\n- Holiday Name`
   - **Leave**: `วันนี้ใครลา : (2025-08-12)\n- Employee Name`

## API Endpoints

The application currently runs as a command-line tool. Future versions may include REST API endpoints.

## Dependencies

Key dependencies include:
- `google.golang.org/api` - Google Calendar API client
- `github.com/line/line-bot-sdk-go` - Line Bot SDK
- `golang.org/x/oauth2` - OAuth2 authentication

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass and maintain 100% coverage
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Deployment

The project includes GitHub Actions for automated deployment to AWS Lambda. See `.github/workflows/deploy.yml` for deployment configuration.

## Troubleshooting

### Common Issues

1. **Google Calendar API Errors**:
   - Verify service account has access to the calendar
   - Check if Calendar API is enabled in Google Cloud Console
   - Ensure credentials are properly base64 encoded

2. **Line Bot Errors**:
   - Verify bot is added to the target Line group
   - Check channel secret and token are correct
   - Ensure Group ID is valid

3. **Timezone Issues**:
   - The application uses Asia/Bangkok timezone
   - Ensure your system supports this timezone

For more help, please check the GitHub issues or create a new issue with detailed error information.

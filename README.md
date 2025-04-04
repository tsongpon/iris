Iris
===========

![Deploy Lambda](https://github.com/tsongpon/iris/actions/workflows/deploy.yml/badge.svg)

Read event from Google Calendar and notify events to Line group

### Run

To run program from your local machine, please set this following environment variable.

```bash
export GOOGLE_CREDENTIALS_JSON=$(cat your-google-calendar-credential-file | base64)
export CALENDAR_ID=your-calendar-id
export LINE_CHANNEL_SECRET=your-line-channel-secret
export LINE_CHANNEL_TOKEN=your-line-channel-token
export LINE_GROUP_ID=your-line-group-id-to-send-message-to
```

Run program

```bash
go run cmd/main.go
```
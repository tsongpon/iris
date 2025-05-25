package gateway

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineNotificationGateway struct {
	lineGroupID   string
	channelSecret string
	channelToken  string
}

func NewLineNotificationGateway(lineGroupID string, channelSecret string, channelToken string) *LineNotificationGateway {
	return &LineNotificationGateway{lineGroupID: lineGroupID, channelSecret: channelSecret, channelToken: channelToken}
}

func (l *LineNotificationGateway) Send(message string) error {
	lineBot, err := linebot.New(l.channelSecret, l.channelToken)
	if err != nil {
		log.Printf("Failed to create LINE bot: %v", err)
		return err
	}

	log.Printf("Sending message to LINE group")
	_, err = lineBot.PushMessage(l.lineGroupID, linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	return nil
}

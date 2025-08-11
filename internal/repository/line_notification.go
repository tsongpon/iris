package repository

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineNotificationRepository struct {
	lineGroupID   string
	channelSecret string
	channelToken  string
}

func NewLineNotificationRepository(lineGroupID string, channelSecret string, channelToken string) LineNotificationRepository {
	return LineNotificationRepository{lineGroupID: lineGroupID, channelSecret: channelSecret, channelToken: channelToken}
}

func (l LineNotificationRepository) SendNotification(message string) error {
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

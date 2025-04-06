package notichannel

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineNoti struct {
	lineGroupID   string
	channelSecret string
	channelToken  string
}

func NewLineNoti(lineGroupID string, channelSecret string, channelToken string) *LineNoti {
	return &LineNoti{lineGroupID: lineGroupID, channelSecret: channelSecret, channelToken: channelToken}
}

func (l *LineNoti) Send(message string) error {
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

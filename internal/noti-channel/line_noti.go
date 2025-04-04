package notichannel

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func SendLineNoti(message string) error {

	// use group ID C9c2917e798af6f232daba62dfb717b9c for real Open API group
	lineGroupID := os.Getenv("LINE_GROUP_ID")
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")

	lineBot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return err
	}

	log.Printf("Sending message to LINE group")
	_, err = lineBot.PushMessage(lineGroupID, linebot.NewTextMessage(message)).Do()
	if err != nil {
		return err
	}

	return nil
}

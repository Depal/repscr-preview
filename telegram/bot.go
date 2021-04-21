package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Token  string = "[СТЁРТО ОК]"
	ChatID int64  = -1337
)

func SendFileToReportChannel(filePath string, caption string, silent bool) (messageID int, err error) {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		return -1, err
	}

	doc := tgbotapi.NewDocumentUpload(ChatID, filePath)
	doc.Caption = caption
	doc.DisableNotification = silent

	msg, err := bot.Send(doc)
	if err != nil {
		return -1, err
	}

	return msg.MessageID, nil
}

func DeleteMessageFromReportChannel(messageID int) (err error) {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		return err
	}

	deleteMessageConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    ChatID,
		MessageID: messageID,
	}
	_, err = bot.DeleteMessage(deleteMessageConfig)
	if err != nil {
		return err
	}

	return nil
}

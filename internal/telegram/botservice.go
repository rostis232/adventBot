package telegram

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type BotService struct {
	MsgChan chan telego.SendMessageParams
}

func NewBotService (msgChan chan telego.SendMessageParams) *BotService{
	return &BotService{
		MsgChan: msgChan,
	}
}

func (s *BotService) CostructMessage(received *telego.Message) {
	text := "Я можу тебе вислухати, але краще <a href=\"https://www.instagram.com/mariechka232/\">звернись до коуча</a>!\n"
	
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: text,
		ReplyToMessageID: received.MessageID,
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
	
}
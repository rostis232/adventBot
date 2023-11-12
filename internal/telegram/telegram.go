package telegram

import (
	"log"

	"github.com/mymmrac/telego"
)

type Bot struct {
	tgBot *telego.Bot
	msgChan chan telego.SendMessageParams
	botService *BotService
}

func NewBot(token string) *Bot {
	tg := &Bot{}
	bot, err := telego.NewBot(token, telego.WithWarnings())
	if err != nil {
		log.Fatalln(err)
	}
	tg.tgBot = bot
	tg.msgChan = make(chan telego.SendMessageParams)
	tg.botService = NewBotService(tg.msgChan)
	return tg
}

func (b *Bot) ListenTelegram() {
	updates, err := b.tgBot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Println(err)
	}
	
	for update := range updates {
		if update.Message != nil {
			b.botService.CostructMessage(update.Message)
		}
	}
}

func (b *Bot) SendMessages() {
	for msg := range b.msgChan {
		_, err := b.tgBot.SendMessage(&msg)
		if err != nil {
			log.Println(err)
		}
	}
	
}
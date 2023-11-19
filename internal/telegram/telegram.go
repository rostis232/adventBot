package telegram

import (
	"log"

	"github.com/mymmrac/telego"
	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/repository"
)

type Bot struct {
	tgBot *telego.Bot
	MsgChan chan telego.SendMessageParams
	BotService *BotService
}

func NewBot(token string, config *config.Config, repo *repository.Repository) *Bot {
	tg := &Bot{}
	bot, err := telego.NewBot(token, telego.WithWarnings())
	if err != nil {
		log.Fatalln(err)
	}
	tg.tgBot = bot
	tg.MsgChan = make(chan telego.SendMessageParams)
	tg.BotService = NewBotService(tg.MsgChan, config, repo)
	return tg
}

func (b *Bot) ListenTelegram() {
	updates, err := b.tgBot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Println(err)
	}
	
	for update := range updates {
		if update.Message != nil {
			b.BotService.BotRouter(update.Message)
		}
	}
}

func (b *Bot) SendMessages() {
	for msg := range b.MsgChan {
		_, err := b.tgBot.SendMessage(&msg)
		if err != nil {
			log.Println(err)
		}
	}
	
}
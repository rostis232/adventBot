package telegram

import (
	"database/sql"
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/repository"
)

type BotService struct {
	MsgChan chan telego.SendMessageParams
	Config *config.Config
	Repo *repository.Repository
}

func NewBotService (msgChan chan telego.SendMessageParams, config *config.Config, repo *repository.Repository) *BotService{
	return &BotService{
		MsgChan: msgChan,
		Config: config,
		Repo: repo,
	}
}

func (s *BotService) BotRouter(received *telego.Message) {
	fmt.Println(received.Text)
	//Get info about chatID from DB
	costumer, err := s.Repo.GetCostumerByChatID(int(received.Chat.ID))
	switch {
	//If we got an error
	case err != nil && err != sql.ErrNoRows:
		fmt.Println("Error while reading costumer from DB:", err, received)
		s.ErrorMessage(received)
	//Check if chat ID is registered
	case err == sql.ErrNoRows:
		if err := s.Repo.AddCostumer(int(received.Chat.ID)); err != nil {
			fmt.Println("Error while adding new costumer:", err, received)
			s.ErrorMessage(received)
			return
		}
		s.NeedNameMessage(received)
	//If User is registered, check if he has any status:
	//0 - undefined
	//1 - waiting for name
	case *costumer.Status == 1:
		if err := s.Repo.ChangeNameAndStatusTo2(int(received.Chat.ID), received.Text); err != nil {
			fmt.Println("Error while changing name:", err, received)
			s.ErrorMessage(received)
			return
		}
		s.NeedSecretKey(received)
	case *costumer.Status == 2:
		rws, err := s.Repo.SetRelationWithSecretKey(*costumer.CostumerID, received.Text)
		fmt.Println(rws)
		if err != nil {
			fmt.Println("Error while setting relation with secret key:", err, received)
			s.ErrorMessage(received)
			return
		}
		if rws == 0 {
			s.WrongSecretKey(received)
			return
		}
		s.SecretKeyAcepted(received)
		if err := s.Repo.ChangeStatusTo3(int(received.Chat.ID)); err != nil {
			fmt.Println("Error while changing status to 3:", err, received)
			s.ErrorMessage(received)
			return
		}
	default:
		s.UndefinedMessage(received)
	}

	
}

func (s *BotService) UndefinedMessage (received *telego.Message){
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: "Я можу тебе вислухати, але краще <a href=\""+s.Config.InstaLink+"\">звернись до коуча</a>!\n",
		ReplyToMessageID: received.MessageID,
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) ErrorMessage (received *telego.Message){
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: "Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) NeedNameMessage (received *telego.Message){
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: "Круто тебе тут бачити! Я адвент-календар <a href=\""+s.Config.InstaLink+"\">коуча Марії</a>!\nЯк мені до тебе звертатись?\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) NeedSecretKey (received *telego.Message) {
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: received.Text+", радий знайомству! Якщо в тебе вже є Ключ до мене - введи його, якщо ні - звернись до <a href=\""+s.Config.InstaLink+"\">коуча Марії</a>!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) SecretKeyAcepted (received *telego.Message) {
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: "Все вірно! Це маленький крок, який наблизив тебе до розуміння себе! Очікуй від мене щоденних повідомлень і обов'язково виконуй мої інструкції! В подарунок отримаєш.... себе!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) WrongSecretKey (received *telego.Message) {
	msg := telego.SendMessageParams{
		ChatID: tu.ID(received.Chat.ID),
		Text: "Я подивився свої записи, але не знайшов такого ключа. Його точно правильно введено? Звернись до <a href=\""+s.Config.InstaLink+"\">коуча Марії</a>!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}
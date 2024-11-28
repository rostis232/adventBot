package telegram

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/models"
	"github.com/rostis232/adventBot/internal/repository"
)

type BotService struct {
	MsgChan chan telego.SendMessageParams
	Config  *config.Config
	Repo    *repository.Repository
}

func NewBotService(msgChan chan telego.SendMessageParams, config *config.Config, repo *repository.Repository) *BotService {
	return &BotService{
		MsgChan: msgChan,
		Config:  config,
		Repo:    repo,
	}
}

func (s *BotService) BotRouter(received *telego.Message) {
	//TODO:ADD link reading
	fmt.Println(received.Text)
	//Get info about costumer with a chatID from DB
	costumer, err := s.Repo.GetCostumerByChatID(int(received.Chat.ID))
	switch {

	//If we got an error
	case err != nil && err != sql.ErrNoRows:
		fmt.Println("Error while reading costumer from DB:", err, received)
		s.ErrorMessage(received)

	//if recievded registration link
	case strings.Contains(received.Text, "/start "):
		s.AceptedRegistrationLink(costumer, received)

	//If user is unregistered
	case err == sql.ErrNoRows:
		s.NewCostumer(received)

	//If we expecting code
	case *costumer.WaitingFor == repository.WaitingForCode:
		s.ExpectingCode(costumer, received)

	//If we expecting name
	case *costumer.WaitingFor == repository.WaitingForName:
		s.ExpectingName(costumer, received)

	default:
		s.UndefinedMessage(received)
	}
}

func (s *BotService) AceptedRegistrationLink(costumer models.Costumer, received *telego.Message) {
	var err error
	code, _ := strings.CutPrefix(received.Text, "/start ")
	if costumer.ChatID == nil || *costumer.ChatID == 0 {
		if err := s.Repo.AddCostumer(int(received.Chat.ID)); err != nil {
			fmt.Println("Error while adding new costumer:", err, received)
			s.ErrorMessage(received)
			return
		}
		costumer, err = s.Repo.GetCostumerByChatID(int(received.Chat.ID))
		if err != nil {
			fmt.Println("Error while getting info about just registered costumer:", err, received)
			s.ErrorMessage(received)
			return
		}
	}

	if costumer.IsActivated != nil && *costumer.IsActivated == 1 {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Твій обліковий запис вже активовано!",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}

	rws, err := s.Repo.SetRelationWithSecretKey(*costumer.ChatID, code)
	if err != nil {
		fmt.Println("Error while setting relation with secret key:", err, received)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}
	if rws == 0 {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Я подивився свої записи, але не знайшов такого ключа. Його точно правильно введено? Якщо була помилка відправ його мені ще раз, або звернись до <a href=\"" + s.Config.InstaLink + "\">коуча Марії</a>!\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err := s.Repo.SetStatusWaitingForCode(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
		return
	}

	err = s.Repo.SetActivated(int(received.Chat.ID))

	if err != nil {
		fmt.Println(err)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}

	if costumer.Name == nil || *costumer.Name == "" {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Круто тебе тут бачити! Я адвент-календар <a href=\"" + s.Config.InstaLink + "\">коуча Марії</a>!\nЯк мені до тебе звертатись?\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForName(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	} else {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Круто тебе тут бачити! Я адвент-календар <a href=\"" + s.Config.InstaLink + "\">коуча Марії</a>!\n Це маленький крок, який наблизив тебе до розуміння себе! Очікуй від мене щоденних повідомлень і обов'язково виконуй мої інструкції! В подарунок отримаєш.... себе!\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForNothing(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	}
}

func (s *BotService) UndefinedMessage(received *telego.Message) {
	msg := telego.SendMessageParams{
		ChatID:           tu.ID(received.Chat.ID),
		Text:             "Зараз я не очікую від тебе повідомлень. Можу тебе вислухати, але краще <a href=\"" + s.Config.InstaLink + "\">звернись до коуча</a>!\n",
		ReplyToMessageID: received.MessageID,
		ParseMode:        "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) ErrorMessage(received *telego.Message) {
	msg := telego.SendMessageParams{
		ChatID:    tu.ID(received.Chat.ID),
		Text:      "Ой, щось пішло не так, як я хотів. Скажи про це <a href=\"" + s.Config.InstaLink + "\">Марії</a>!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) NewCostumer(received *telego.Message) {
	if err := s.Repo.AddCostumer(int(received.Chat.ID)); err != nil {
		fmt.Println("Error while adding new costumer:", err, received)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}
	msg := telego.SendMessageParams{
		ChatID:    tu.ID(received.Chat.ID),
		Text:      "Радий тебе тут бачити! Якщо в тебе вже є Ключ до мене - введи його, якщо ні - звернись до <a href=\"" + s.Config.InstaLink + "\">коуча Марії</a>!\n",
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
	err := s.Repo.SetStatusWaitingForCode(int(received.Chat.ID))
	if err != nil {
		fmt.Println(err)
		msg = telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}
}

func (s *BotService) ExpectingCode(costumer models.Costumer, received *telego.Message) {
	rws, err := s.Repo.SetRelationWithSecretKey(int(received.Chat.ID), received.Text)
	if err != nil {
		fmt.Println("Error while setting relation with secret key:", err, received)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}
	if rws == 0 {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Я подивився свої записи, але не знайшов такого ключа. Або він вже зайнятий. Його точно правильно введено? Якщо була помилка відправ його мені ще раз, або звернись до <a href=\"" + s.Config.InstaLink + "\">коуча Марії</a>!\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		if *costumer.WaitingFor != repository.WaitingForCode {
			err := s.Repo.SetStatusWaitingForCode(int(received.Chat.ID))
			if err != nil {
				fmt.Println(err)
				msg = telego.SendMessageParams{
					ChatID:    tu.ID(received.Chat.ID),
					Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
					ParseMode: "HTML",
				}
				s.MsgChan <- msg
				return
			}
		}
		return
	}

	err = s.Repo.SetActivated(int(received.Chat.ID))
	if err != nil {
		fmt.Println(err)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}

	if costumer.Name == nil || *costumer.Name == "" {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Все вірно! Тепер давай знайомитись! Як тебе звуть?\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForName(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	} else {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Радий нашому знайомству! 🎉\nПопереду на вас чекає 25 завдань, які допоможуть краще зрозуміти свої фінанси, зробити маленькі кроки до великих цілей і відкрити нові можливості. 💡\n\nКожен день – це новий виклик, але й новий шанс змінити своє життя! 🚀\nГотові розпочати цю захопливу подорож? 🌟",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForNothing(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	}
}

func (s *BotService) ExpectingName(costumer models.Costumer, received *telego.Message) {
	err := s.Repo.ChangeName(int(received.Chat.ID), received.Text)
	if err != nil {
		fmt.Println("Error while setting relation with secret key:", err, received)
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		return
	}

	if costumer.IsActivated == nil || *costumer.IsActivated == 0 {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Радий знайомству.\nЩоб продовжити мені потрібно отримати код. Отримати його можна у <a href=\"" + s.Config.InstaLink + "\">Марії</a>.\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForCode(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	} else {
		msg := telego.SendMessageParams{
			ChatID:    tu.ID(received.Chat.ID),
			Text:      "Радий знайомству! Це маленький крок, який наблизив тебе до розуміння себе! Очікуй від мене щоденних повідомлень і обов'язково виконуй мої інструкції! В подарунок отримаєш.... себе!\n",
			ParseMode: "HTML",
		}
		s.MsgChan <- msg
		err = s.Repo.SetStatusWaitingForNothing(int(received.Chat.ID))
		if err != nil {
			fmt.Println(err)
			msg = telego.SendMessageParams{
				ChatID:    tu.ID(received.Chat.ID),
				Text:      fmt.Sprint("Ой, щось пішло не так, як я хотів. Скажи про це <a href=\""+s.Config.InstaLink+"\">Марії</a>: %s\n", err),
				ParseMode: "HTML",
			}
			s.MsgChan <- msg
			return
		}
	}
}

func (s *BotService) SendMessageNow(chatID int, message string) {
	msg := telego.SendMessageParams{
		ChatID:    tu.ID(int64(chatID)),
		Text:      message,
		ParseMode: "HTML",
	}
	s.MsgChan <- msg
}

func (s *BotService) CheckUnsendedMessages() {
	messages, err := s.Repo.GetAllUnsendedMessages()
	if err != nil {
		fmt.Println(err)
		return
	}
	costumers, err := s.Repo.GetAllActivatedCustumers()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range messages {
		dateTime, err := time.Parse("2006-01-02 15:04 -0700", v.DateTime)
		if err != nil {
			fmt.Println(err)
			return
		}

		if dateTime.Compare(time.Now()) < 0 {
			for _, c := range costumers {
				fmt.Println("Відправляю", v.MessageID)
				s.SendMessageNow(*c.ChatID, v.Text)
			}

			err = s.Repo.SetStatusSent(v.MessageID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}

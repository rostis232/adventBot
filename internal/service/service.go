package service

import (
	"fmt"
	"time"

	"github.com/rostis232/adventBot/internal/models"
	"github.com/rostis232/adventBot/internal/telegram"
)

type Repository interface{
	GetAllCustumers() ([]models.Costumer, error)
	GetAllActivatedCustumers() ([]models.Costumer, error)
	GetCostumerByChatID(chatID int) (models.Costumer, error)
	AddCostumer(chatID int) error
	ChangeName(chatID int, name string) error
	SetRelationWithSecretKey(costumerID int, secretKey string) (int, error)
	GetAllMessages() ([]models.Message, error)
	AddMessage(dateTime, message string) error
	GetAllUnsendedMessages() ([]models.Message, error)
	GetAllSecretKeys() ([]models.SecretKey, error)
	AddKey(key int) error
}

type Service struct{
	Repo Repository
	bot *telegram.Bot
}

func NewService (repo Repository, bot *telegram.Bot) *Service{
	return &Service{Repo: repo,
	bot: bot,}
}

func (s *Service) SendMessageNow (message string) error {
	// get all active users
	costumers, err := s.Repo.GetAllActivatedCustumers()
	if err != nil {
		return err
	}

	// send to all of them a message
	for _, c := range costumers {
		s.bot.BotService.SendMessageNow(*c.ChatID, message)
	}

	return nil
}

func (s *Service) GetAllMessages() ([]models.Message, error) {
	return s.Repo.GetAllMessages()
}

func (s *Service) AddMessage(dateTime, message string) error {
	return s.Repo.AddMessage(dateTime, message)
}

func (s *Service) GetAllSecretKeys() ([]models.SecretKey, error) {
	return s.Repo.GetAllSecretKeys()
}

func (s *Service) AddNewKeys() error {
	now := time.Now()
	fmt.Println(now.Unix())
	return s.Repo.AddKey(int(now.Unix()))
}

func (s *Service) GetAllCustumers() ([]models.Costumer, error) {
	return s.Repo.GetAllCustumers()
}
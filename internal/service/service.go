package service

import (
	"github.com/rostis232/adventBot/internal/models"
	"github.com/rostis232/adventBot/internal/telegram"
)

type Repository interface{
	GetAllCustumers() ([]models.Costumer, error)
	GetCostumerByChatID(chatID int) (models.Costumer, error)
	AddCostumer(chatID int) error
	ChangeNameAndStatusTo2(chatID int, name string) error
	SetRelationWithSecretKey(costumerID int, secretKey string) (int, error)
	ChangeStatusTo3(chatID int) error
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
	costumers, err := s.Repo.GetAllCustumers()
	if err != nil {
		return err
	}

		// send to all of them a message
	for _, c := range costumers {
		if *c.Status != 3 {
			continue
		} else {
			s.bot.BotService.SendMessageNow(*c.ChatID, message)
		}
	}

	return nil
}

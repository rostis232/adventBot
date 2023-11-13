package repository

import (
	"database/sql"

	"github.com/rostis232/adventBot/internal/models"
)

type Repository struct{
	DB *sql.DB
}

func NewRepository (db *sql.DB) *Repository{
	return &Repository{DB: db}
}

func(r *Repository) GetAllCustumers(){}

func(r *Repository) GetCostumerByChatID(chatID int) (models.Costumer, error) {
	costumer := models.Costumer{}

	query := "SELECT costumer_id, chat_id, name, status FROM costumers WHERE chat_id = ?"
	row := r.DB.QueryRow(query, chatID)

	err := row.Scan(&costumer.CostumerID, &costumer.ChatID, &costumer.Name, &costumer.Status)
	if err != nil {
		return costumer, err
	}
	return costumer, nil
}

func (r *Repository) AddCostumer(chatID int) error {
	query := "INSERT INTO costumers (chat_id, status) VALUES (?,1)"
	_, err := r.DB.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangeNameAndStatusTo2(chatID int, name string) error {
	query := "UPDATE costumers SET name = ?, status = 2 WHERE chat_id = ?"
	_, err := r.DB.Exec(query, name, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetRelationWithSecretKey(costumerID int, secretKey string) (int, error) {
	query := "UPDATE secret_keys SET costumer_id = ? WHERE secret_key = ? AND costumer_id IS NULL"
	result, err := r.DB.Exec(query, costumerID, secretKey)
	if err != nil {
		return 0, err
	}
	rws, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rws), nil
}

func (r *Repository) ChangeStatusTo3(chatID int) error {
	query := "UPDATE costumers SET status = 3 WHERE chat_id = ?"
	_, err := r.DB.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}
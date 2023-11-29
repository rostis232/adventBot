package repository

import (
	"database/sql"
	"fmt"

	"github.com/rostis232/adventBot/internal/models"
)

type Repository struct{
	DB *sql.DB
}

func NewRepository (db *sql.DB) *Repository{
	return &Repository{DB: db}
}

var (
	WaitingForNothing = 0
	WaitingForName = 1
	WaitingForCode = 2
)

func(r *Repository) GetAllCustumers() ([]models.Costumer, error) {
	query := "SELECT chat_id, name, is_activated, waiting_for FROM costumers"
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.Costumer{}, err
	}
	defer rows.Close()

	var costumers []models.Costumer

	// Читання результатів запиту
	for rows.Next() {
		var costumer models.Costumer
		if err := rows.Scan(&costumer.ChatID, &costumer.Name, &costumer.IsActivated, &costumer.WaitingFor); err != nil {
			return costumers, err
		}
		costumers = append(costumers, costumer)
	}

	if err := rows.Err(); err != nil {
		return costumers, err
	}

	return costumers, nil
}

func(r *Repository) GetAllActivatedCustumers() ([]models.Costumer, error) {
	query := "SELECT chat_id, name, is_activated, waiting_for FROM costumers WHERE is_activated = 1"
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.Costumer{}, err
	}
	defer rows.Close()

	var costumers []models.Costumer

	// Читання результатів запиту
	for rows.Next() {
		var costumer models.Costumer
		if err := rows.Scan(&costumer.ChatID, &costumer.Name, &costumer.IsActivated, &costumer.WaitingFor); err != nil {
			return costumers, err
		}
		costumers = append(costumers, costumer)
	}

	if err := rows.Err(); err != nil {
		return costumers, err
	}

	return costumers, nil
}

func(r *Repository) GetCostumerByChatID(chatID int) (models.Costumer, error) {
	costumer := models.Costumer{}

	query := "SELECT chat_id, name, is_activated, waiting_for FROM costumers WHERE chat_id = ?"
	row := r.DB.QueryRow(query, chatID)

	err := row.Scan(&costumer.ChatID, &costumer.Name, &costumer.IsActivated, &costumer.WaitingFor)
	if err != nil {
		return costumer, err
	}
	return costumer, nil
}

func (r *Repository) AddCostumer(chatID int) error {
	query := "INSERT INTO costumers (chat_id, is_activated, waiting_for) VALUES (?, 0, 0)"
	_, err := r.DB.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangeName(chatID int, name string) error {
	query := "UPDATE costumers SET name = ? WHERE chat_id = ?"
	_, err := r.DB.Exec(query, name, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetRelationWithSecretKey(chatID int, secretKey string) (int, error) {
	query := "UPDATE secret_keys SET chat_id = ? WHERE secret_key = ? AND chat_id IS NULL"
	result, err := r.DB.Exec(query, chatID, secretKey)
	if err != nil {
		return 0, err
	}
	rws, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rws), nil
}

func(r *Repository) GetAllMessages() ([]models.Message, error) {
	query := "SELECT message_id, date, message, is_sent FROM messages"
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.Message{}, err
	}
	defer rows.Close()

	var messages []models.Message

	// Читання результатів запиту
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.MessageID, &message.DateTime, &message.Text, &message.Sent); err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func(r *Repository) AddMessage(dateTime, message string) error {
	query := "INSERT INTO messages (date, message, is_sent) VALUES (?, ?, 0)"
	_, err := r.DB.Exec(query, dateTime, message)
	if err != nil {
		return err
	}
	return nil
}

func(r *Repository) GetAllUnsendedMessages() ([]models.Message, error) {
	query := "SELECT message_id, date, message, is_sent FROM messages WHERE is_sent = 0"
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.Message{}, err
	}
	defer rows.Close()

	var messages []models.Message

	// Читання результатів запиту
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.MessageID, &message.DateTime, &message.Text, &message.Sent); err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func(r *Repository) SetStatusSent(id int) error {
	query := "UPDATE messages SET is_sent = 1 WHERE message_id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func(r *Repository) SetStatusWaitingForNothing(chatID int) error {
	query := "UPDATE costumers SET waiting_for = ? WHERE chat_id = ?"
	_, err := r.DB.Exec(query, WaitingForNothing, chatID)
	if err != nil {
		return err
	}
	return nil
}

func(r *Repository) SetStatusWaitingForName(chatID int) error {
	query := "UPDATE costumers SET waiting_for = ? WHERE chat_id = ?"
	_, err := r.DB.Exec(query, WaitingForName, chatID)
	if err != nil {
		return err
	}
	return nil
}

func(r *Repository) SetStatusWaitingForCode(chatID int) error {
	query := "UPDATE costumers SET waiting_for = ? WHERE chat_id = ?"
	_, err := r.DB.Exec(query, WaitingForCode, chatID)
	if err != nil {
		return err
	}
	return nil
}

func(r *Repository) SetActivated(chatID int) error {
	fmt.Println(chatID)
	query := "UPDATE costumers SET is_activated = 1 WHERE chat_id = ?"
	result, err := r.DB.Exec(query, chatID)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func(r *Repository) GetAllSecretKeys() ([]models.SecretKey, error) {
	query := "SELECT sk_id, secret_key, chat_id FROM secret_keys"
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.SecretKey{}, err
	}
	defer rows.Close()

	var keys []models.SecretKey

	// Читання результатів запиту
	for rows.Next() {
		var key models.SecretKey
		if err := rows.Scan(&key.SkID, &key.SecretKey, &key.ChatID); err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return keys, err
	}

	return keys, nil
}
package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//write here constants for db entities

func NewSQLiteDB (file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./"+file)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS "messages" (
		"message_id"	INTEGER UNIQUE,
		"date"	TEXT NOT NULL,
		"message"	TEXT NOT NULL,
		"is_sent"	INTEGER DEFAULT 0,
		PRIMARY KEY("message_id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "costumers" (
		"costumer_id"	INTEGER UNIQUE,
		"chat_id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT,
		"status" INTEGER,
		PRIMARY KEY("costumer_id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "secret_keys" (
		"sk_id"	INTEGER UNIQUE,
		"secret_key"	INTEGER NOT NULL UNIQUE,
		"costumer_id"	INTEGER DEFAULT NULL UNIQUE,
		FOREIGN KEY("costumer_id") REFERENCES "costumers"("costumer_id"),
		PRIMARY KEY("sk_id" AUTOINCREMENT)
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}
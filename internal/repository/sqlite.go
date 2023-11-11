package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//write here constants for db entities

func NewSQLiteDB (file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
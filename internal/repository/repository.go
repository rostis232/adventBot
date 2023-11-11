package repository

import "database/sql"

type Repository struct{
	DB *sql.DB
}

func NewRepository (db *sql.DB) *Repository{
	return &Repository{DB: db}
}

func(r *Repository) GetAllCustumers(){}
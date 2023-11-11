package app

import (
	"github.com/rostis232/adventBot/internal/handler"
	"github.com/rostis232/adventBot/internal/repository"
	"github.com/rostis232/adventBot/internal/service"
)

type App struct {
	repo *repository.Repository
	service *service.Service
	handler *handler.Handler
}

func NewApp (dbName string) (*App, error) {
	app := &App{}
	db, err := repository.NewSQLiteDB(dbName)
	if err != nil {
		return nil, err
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	app.repo = repo
	app.service = service
	app.handler = handler

	return app, nil
}

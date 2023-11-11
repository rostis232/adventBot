package app

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"github.com/rostis232/adventBot/internal/handler"
	"github.com/rostis232/adventBot/internal/repository"
	"github.com/rostis232/adventBot/internal/service"
)

type App struct {
	repo *repository.Repository
	service *service.Service
	handler *handler.Handler
	echo *echo.Echo
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

	//ECHO
	app.echo = echo.New()
	app.echo.Use(middleware.Logger())
	app.echo.Use(middleware.Recover())

	//Make Routes

	return app, nil
}

func (a *App) Run(port string) {
	// Start server
	a.echo.Logger.Fatal(a.echo.Start(":"+port))
}

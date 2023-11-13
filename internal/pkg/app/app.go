package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/handler"
	"github.com/rostis232/adventBot/internal/renderer"
	"github.com/rostis232/adventBot/internal/repository"
	"github.com/rostis232/adventBot/internal/service"
	"github.com/rostis232/adventBot/internal/telegram"
)

type App struct {
	config *config.Config
	repo *repository.Repository
	service *service.Service
	handler *handler.Handler
	echo *echo.Echo
	bot *telegram.Bot
}

func NewApp (config *config.Config) (*App, error) {
	app := &App{}
	app.config = config
	db, err := repository.NewSQLiteDB(app.config.DBname)
	if err != nil {
		return nil, err
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(app.config, service)

	app.repo = repo
	app.service = service
	app.handler = handler

	//ECHO
	app.echo = echo.New()
	app.echo.Use(middleware.Logger())
	app.echo.Use(middleware.Recover())
	app.echo.Renderer = renderer.Tmps

	//Make Routes
	app.echo.GET("/", app.handler.Home)

	bot := telegram.NewBot(config.TGsecretCode, app.config, app.repo)
	app.bot = bot

	return app, nil
}

func (a *App) Run() {
	// Start server
	go a.bot.ListenTelegram()
	go a.bot.SendMessages()
	a.echo.Logger.Fatal(a.echo.Start(":"+a.config.Port))
}

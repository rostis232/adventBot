package app

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/handler"
	"github.com/rostis232/adventBot/internal/renderer"
	"github.com/rostis232/adventBot/internal/repository"
	"github.com/rostis232/adventBot/internal/service"
	"github.com/rostis232/adventBot/internal/sessions"
	"github.com/rostis232/adventBot/internal/telegram"
	session "github.com/spazzymoto/echo-scs-session"
)

type App struct {
	config *config.Config
	repo *repository.Repository
	service *service.Service
	handler *handler.Handler
	echo *echo.Echo
	bot *telegram.Bot
	session *sessions.Sessions
}

func NewApp (config *config.Config) (*App, error) {
	app := &App{}
	app.config = config
	
	db, err := repository.NewSQLiteDB(app.config.DBname)
	if err != nil {
		return nil, err
	}
	sess := sessions.NewSessions(app.config.RedisAddress)
	repo := repository.NewRepository(db)
	bot := telegram.NewBot(config.TGsecretCode, app.config, repo)
	app.bot = bot
	srvs := service.NewService(repo, bot)
	hdl := handler.NewHandler(app.config, srvs, sess)

	app.repo = repo
	app.service = srvs
	app.handler = hdl
	app.session = sess

	//ECHO
	app.echo = echo.New()
	app.echo.Use(middleware.Logger())
	app.echo.Use(middleware.Recover())
	app.echo.Renderer = renderer.Tmps

	//Make Routes
	app.echo.Use(session.LoadAndSave(app.session.SessionManager))
	app.echo.GET("/", app.handler.Home)
	app.echo.GET("/login", app.handler.Login)
	app.echo.POST("/login", app.handler.PostLoginPage)
	app.echo.GET("/logout", app.handler.Logout)
	app.echo.GET("/send-page", app.handler.SendPage)
	app.echo.POST("/send-page", app.handler.PostSendPage)
	app.echo.GET("/journal", app.handler.Journal)
	app.echo.GET("/journal/add", app.handler.JournalAdd)
	app.echo.POST("/journal/add", app.handler.PostJournalAdd)
	app.echo.GET("/keys", app.handler.Keys)

	

	return app, nil
}

func (a *App) Run() {
	// Start server
	go a.bot.ListenTelegram()
	go a.bot.SendMessages()
	go a.bot.CheckUnsendedMessages(5 * time.Second)
	a.echo.Logger.Fatal(a.echo.Start(":"+a.config.Port))
}
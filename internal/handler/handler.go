package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rostis232/adventBot/internal/models"
	"github.com/rostis232/adventBot/internal/sessions"

	"github.com/labstack/echo/v4"
	"github.com/rostis232/adventBot/config"
	templatedata "github.com/rostis232/adventBot/internal/template_data"
)

type Service interface{
	SendMessageNow (message string) error
	GetAllMessages() ([]models.Message, error)
	AddMessage(dateTime, message string) error
}

type Handler struct {
	config  *config.Config
	Service Service
	Session *sessions.Sessions
}

func NewHandler(config *config.Config, service Service, session *sessions.Sessions) *Handler {
	return &Handler{config: config,
		Service: service,
		Session: session,
	}
}

func (h *Handler) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", h.AddDefaultData(nil, c.Request()))
}

func (h *Handler) Login (c echo.Context) error {
	return c.Render(http.StatusOK, "login.page.html", h.AddDefaultData(nil, c.Request()))
}

func (h *Handler) PostLoginPage (c echo.Context) error {
	_ = h.Session.SessionManager.RenewToken(c.Request().Context())
	login := c.FormValue("login")
	pass := c.FormValue("password")
	if login != h.config.AdminLogin && pass != h.config.AdminPass {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Неправильні дані.")
		return c.Redirect(http.StatusSeeOther, "/login")
	} else {
		h.Session.SessionManager.Put(c.Request().Context(), "user", "1")

		h.Session.SessionManager.Put(c.Request().Context(), "flash", "Успішний вхід!")

		// redirect the user
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

func (h *Handler) Logout (c echo.Context) error {
	_ = h.Session.SessionManager.Destroy(c.Request().Context())
	_ = h.Session.SessionManager.RenewToken(c.Request().Context())
	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) SendPage (c echo.Context) error {
	if h.Session.SessionManager.GetString(c.Request().Context(), "user") != "1" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Потрібна авторизація!")
		return c.Redirect(http.StatusSeeOther, "/login")
	} else {
		return c.Render(http.StatusOK, "send.page.html", h.AddDefaultData(nil, c.Request())) 
	}
}

func (h *Handler) PostSendPage (c echo.Context) error {
	if h.Session.SessionManager.GetString(c.Request().Context(), "user") != "1" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Потрібна авторизація!")
		return c.Redirect(http.StatusSeeOther, "/login")
	} else {
		message := c.FormValue("message")
		if message == "" {
			h.Session.SessionManager.Put(c.Request().Context(), "warning", "Повідомлення пусте.")
			return c.Render(http.StatusOK, "send.page.html", h.AddDefaultData(nil, c.Request())) 
		} else {
			err := h.Service.SendMessageNow(message)
			if err != nil {
				h.Session.SessionManager.Put(c.Request().Context(), "error", fmt.Sprintf("Помилка надсилання: %s", err))
				return c.Render(http.StatusOK, "send.page.html", h.AddDefaultData(nil, c.Request())) 
			}
			h.Session.SessionManager.Put(c.Request().Context(), "flash", "Повідомлення відправлено.")
			return c.Render(http.StatusOK, "send.page.html", h.AddDefaultData(nil, c.Request())) 
		}
	}
}

func(h *Handler) Journal (c echo.Context) error {
	if h.Session.SessionManager.GetString(c.Request().Context(), "user") != "1" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Потрібна авторизація!")
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	messages, err := h.Service.GetAllMessages()
	if err != nil {
		h.Session.SessionManager.Put(c.Request().Context(), "error", fmt.Sprintf("Помилка отримання переліку повідомлен: %s", err))
		return c.Render(http.StatusOK, "journal.page.html", h.AddDefaultData(nil, c.Request())) 
	}
	td := templatedata.TemplateData{
		Data: map[string]interface{}{"messages":messages},
	}
	return c.Render(http.StatusOK, "journal.page.html", h.AddDefaultData(&td, c.Request())) 
}

func(h *Handler) JournalAdd (c echo.Context) error {
	if h.Session.SessionManager.GetString(c.Request().Context(), "user") != "1" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Потрібна авторизація!")
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	return c.Render(http.StatusOK, "journal_add.page.html", h.AddDefaultData(&templatedata.TemplateData{}, c.Request()))
}

func(h *Handler) PostJournalAdd (c echo.Context) error {
	if h.Session.SessionManager.GetString(c.Request().Context(), "user") != "1" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Потрібна авторизація!")
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	date := c.FormValue("date")
	hour := c.FormValue("hour")
	minutes := c.FormValue("minutes")
	message := c.FormValue("message")
	if date == "" || hour == "" || minutes == "" || message == "" {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Всі поля обов'язкові!")
		return c.Redirect(http.StatusSeeOther, "/journal/add")
	}
	if err := h.Service.AddMessage(date+" "+hour+":"+minutes+" "+h.config.MyTime, message); err != nil {
		h.Session.SessionManager.Put(c.Request().Context(), "error", err)
		return c.Redirect(http.StatusSeeOther, "/journal/add")
	}
	h.Session.SessionManager.Put(c.Request().Context(), "flash", "Нове повідомлення успішно додано!")
	return  c.Redirect(http.StatusSeeOther, "/journal")
}

func(h *Handler) AddDefaultData(td *templatedata.TemplateData, r *http.Request) *templatedata.TemplateData {
	if td == nil {
		td = &templatedata.TemplateData{}
	}
	td.Config = *h.config
	td.Flash = h.Session.SessionManager.PopString(r.Context(), "flash")
	td.Warning = h.Session.SessionManager.PopString(r.Context(), "warning")
	td.Error = h.Session.SessionManager.PopString(r.Context(), "error")
	if h.IsAuthenticated(r) {
		td.Authenticated = true
	}
	td.Now = time.Now()

	return td
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	return h.Session.SessionManager.Exists(r.Context(), "user")
}
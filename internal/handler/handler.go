package handler

import (
	"github.com/rostis232/adventBot/internal/sessions"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rostis232/adventBot/config"
	templatedata "github.com/rostis232/adventBot/internal/template_data"
)

type Service interface{}

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
	data := templatedata.TemplateData{}
	data.Config = *h.config 
	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) Login (c echo.Context) error {
	return c.Render(http.StatusOK, "login.page.html", h.AddDefaultData(nil, c.Request()))
}

func (h *Handler) PostLoginPage (c echo.Context) error {
	_ = h.Session.SessionManager.RenewToken(c.Request().Context())
	login := c.FormValue("login")
	pass := c.FormValue("password")
	if login != h.config.AdminLogin && pass != h.config.AdminPass {
		h.Session.SessionManager.Put(c.Request().Context(), "error", "Invalid credentials.")
		return c.Redirect(http.StatusSeeOther, "/login")
	} else {
		h.Session.SessionManager.Put(c.Request().Context(), "user", 1)

		h.Session.SessionManager.Put(c.Request().Context(), "flash", "Successful login!")

		// redirect the user
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

func (h *Handler) Logout (c echo.Context) error {
	_ = h.Session.SessionManager.Destroy(c.Request().Context())
	_ = h.Session.SessionManager.RenewToken(c.Request().Context())
	return c.Redirect(http.StatusSeeOther, "/login")
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
		//TODO - get more user information
	}
	td.Now = time.Now()

	return td
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	return h.Session.SessionManager.Exists(r.Context(), "user")
}
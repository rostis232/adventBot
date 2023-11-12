package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rostis232/adventBot/config"
	templatedata "github.com/rostis232/adventBot/internal/template_data"
)

type Service interface{}

type Handler struct {
	config  *config.Config
	Service Service
}

func NewHandler(config *config.Config, service Service) *Handler {
	return &Handler{config: config,
		Service: service}
}

func (h *Handler) Home(c echo.Context) error {
	data := templatedata.TemplateData{}
	data.Config = *h.config 
	return c.Render(http.StatusOK, "index.html", data)
}

package transport

import (
	"webpush/config"
	"webpush/pages"

	"github.com/labstack/echo/v4"
)

type PageHandlerInterface interface {
	GetLandingPage(c echo.Context) error
}

type PageHandler struct{ configs config.Configs }

func NewPageHandler(configs config.Configs) PageHandlerInterface {
	return &PageHandler{configs: configs}
}

func (h *PageHandler) GetLandingPage(c echo.Context) error {
	return c.File(pages.LandingPage())
}

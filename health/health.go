package health

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Health struct {
	addr string
	e    *echo.Echo
}

func NewHealth(addr string) (*Health, error) {
	h := &Health{
		addr: addr,
		e:    echo.New(),
	}

	h.e.HideBanner = true
	h.e.HidePort = true
	h.e.DisableHTTP2 = true

	h.e.GET("/healthz", h.health)

	return h, nil
}

func (h *Health) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func (h *Health) Start() error {
	return h.e.Start(h.addr)
}

func (h *Health) Stop() error {
	return h.e.Shutdown(context.Background())
}

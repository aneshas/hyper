package hyperui

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func UseDefaults(e *echo.Echo) error {
	e.HTTPErrorHandler = NewErrorHandler()

	renderer, err := NewRenderer("views/")
	if err != nil {
		return err
	}

	e.Renderer = renderer

	e.Static("/public", "public")
	e.File("/favicon.svg", "public/favicon.svg")

	return nil
}

func Redirect(c echo.Context, path string) error {
	c.Response().Header().Add("Hx-Push-Url", path)

	return c.Redirect(http.StatusFound, path)
}

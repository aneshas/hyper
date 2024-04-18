package hyperui

import "github.com/labstack/echo/v4"

func Defaults(e *echo.Echo) error {
	e.HTTPErrorHandler = NewErrorHandler()

	renderer, err := NewRenderer()
	if err != nil {
		return err
	}

	e.Renderer = renderer

	e.Static("/public", "public")
	e.File("/favicon.svg", "public/favicon.svg")

	return nil
}

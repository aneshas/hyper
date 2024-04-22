package hyperui

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewErrorHandler creates custom http error handler
func NewErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var appErr Error

		if errors.As(err, &appErr) {
			_ = c.Render(appErr.code, appErr.tpl, appErr.Bag)

			return
		}

		_ = c.Render(http.StatusInternalServerError, "server_error", NewBag().WithError("ServerError", err.Error()))
	}
}

func E() Error {
	return Error{
		code: http.StatusBadRequest,
		tpl:  "server_error",
	}
}

type Error struct {
	Bag

	code int
	err  error
	tpl  string
}

func (ve Error) Error() string {
	if ve.err != nil {
		return ve.err.Error()
	}

	return "unknown error"
}

func (ve Error) WithCode(code int) Error {
	ve.code = code

	return ve
}

func (ve Error) WithTpl(name string) Error {
	ve.tpl = name

	return ve
}

func (ve Error) WithBag(bag Bag) Error {
	ve.Bag = bag

	return ve
}

func (ve Error) WithErr(err error) Error {
	ve.err = err

	return ve
}

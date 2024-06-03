package flash

import (
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"time"
)

const (
	flashSuccess = "hyper_flash_success"
	flashWarning = "hyper_flash_warning"
	flashErr     = "hyper_flash_error"
)

func NewMessage(msg string, kv ...string) Message {
	m := Message{
		Message: msg,
		Params:  make(map[string]string),
	}

	for i := 0; i < len(kv); i += 2 {
		m.Params[kv[i]] = kv[i+1]
	}

	return m
}

type Message struct {
	Message string
	Params  map[string]string
}

func (m Message) Encode() string {
	data, _ := json.Marshal(m)

	return encode(data)
}

func Success(c echo.Context, msg string, kv ...string) {
	set(c, flashSuccess, msg, kv...)
}

func Warning(c echo.Context, msg string, kv ...string) {
	set(c, flashWarning, msg, kv...)
}

func Error(c echo.Context, msg string, kv ...string) {
	set(c, flashErr, msg, kv...)
}

func set(c echo.Context, name, msg string, kv ...string) {
	m := NewMessage(msg, kv...)

	c.SetCookie(&http.Cookie{
		Name:  name,
		Value: m.Encode(),
		Path:  "/",
	})
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"flash_success": GetSuccess,
		"flash_warning": GetWarning,
		"flash_error":   GetError,
	}
}

func GetSuccess(c echo.Context) *Message {
	return get(c, flashSuccess)
}

func GetWarning(c echo.Context) *Message {
	return get(c, flashWarning)
}

func GetError(c echo.Context) *Message {
	return get(c, flashErr)
}

func get(c echo.Context, name string) *Message {
	cookie, err := c.Cookie(name)
	if err != nil {
		return nil
	}

	data, err := decode(cookie.Value)
	if err != nil {
		return nil
	}

	var m Message

	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil
	}

	c.SetCookie(
		&http.Cookie{
			Name:    name,
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
			Value:   "",
			Path:    "/",
		},
	)

	return &m
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

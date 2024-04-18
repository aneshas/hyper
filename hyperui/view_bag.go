package hyperui

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

func NewBag() Bag {
	return Bag{
		Data:   make(map[string]any),
		Errors: make(map[string]any),
	}
}

type Bag struct {
	echo.Context

	Data   map[string]any
	Errors map[string]any
}

func (vd Bag) With(val any) Bag {
	return vd.WithKV(reflect.TypeOf(val).Name(), val)
}

func (vd Bag) WithKV(key string, val any) Bag {
	vd.Data[key] = val

	return vd
}

func (vd Bag) WithError(key string, val any) Bag {
	vd.Errors[key] = val

	return vd
}

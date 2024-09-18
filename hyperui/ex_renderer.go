package hyperui

import (
	"fmt"
	"github.com/aneshas/hyper/hyperui/flash"
	"io"

	"github.com/Masterminds/sprig/v3"
	"github.com/dannyvankooten/extemplate"
	"github.com/labstack/echo/v4"
)

func NewRenderer(tplPath string, opts ...Option) (*Renderer, error) {
	xt := extemplate.New().
		Funcs(sprig.FuncMap()).
		Funcs(flash.FuncMap())

	err := xt.ParseDir(tplPath, []string{".tpl"})
	if err != nil {
		return nil, err
	}

	r := &Renderer{
		tpl: xt,
		kv:  make(map[string]any),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

type Option func(r *Renderer)

func WithKV(key string, val any) Option {
	return func(r *Renderer) {
		r.kv[key] = val
	}
}

type Renderer struct {
	tpl *extemplate.Extemplate

	kv map[string]any
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	bag, ok := data.(Bag)
	if ok {
		bag.Context = c

		data = bag.WithKV("Page", name)

		for k, v := range r.kv {
			data = bag.WithKV(k, v)
		}
	}

	return r.tpl.ExecuteTemplate(w, fmt.Sprintf("%s.go.tpl", name), data)
}

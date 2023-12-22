package hyper

import (
	"context"
)

// UCDetails represents ...
type UCDetails struct {
	NS     string
	UCName string
	Req    string
	Resp   string
}

// GeneratedUC represents ...
type GeneratedUC struct {
	Bar string
}

// GenUCFunc accomplishes ...
type GenUCFunc func(ctx context.Context, cmd UCDetails) (*GeneratedUC, error)

type UCGenerator interface {
	GenUC(ctx context.Context, app App, uc UC) error
}

type StoreGenerator interface {
	GenStore(ctx context.Context, app App, uc UC) error
}

// NewGenUC instantiates GenUC use case
func NewGenUC(store AppStore, ucGen UCGenerator, storeGen StoreGenerator) GenUCFunc {
	return func(ctx context.Context, cmd UCDetails) (*GeneratedUC, error) {
		app, err := store.CurrentApp()
		if err != nil {
			return nil, err
		}

		uc := app.NewUC(cmd.NS, cmd.UCName, cmd.Req, cmd.Resp)

		err = ucGen.GenUC(
			ctx, *app, uc,
		)
		if err != nil {
			return nil, err
		}

		return nil, storeGen.GenStore(ctx, *app, uc)
	}
}

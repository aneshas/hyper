package hyperdb

// TODO - Move this to a separate repo?
// implementation should go to hyper boiler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aneshas/tx/boiltx"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Mapper[T any] interface {
	ToEntry(t T) (Entry, error)
	FromEntry(e Entry) (T, error)
}

type Entry interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

type Store[T any, E any] interface {
	ByID(ctx context.Context, id E) (*T, error)
	Save(ctx context.Context, t T) error
}

func NewBoilStore[T any, E any](db *sql.DB, mapper Mapper[T]) *BoilStore[T, E] {
	return &BoilStore[T, E]{
		DB:     db,
		Mapper: mapper,
	}
}

type BoilStore[T any, E any] struct {
	DB     *sql.DB
	Mapper Mapper[T]
}

func (br *BoilStore[T, E]) ByID(ctx context.Context, id E) (*T, error) {
	// TODO implement me
	panic("implement me")
}

func (br *BoilStore[T, E]) Save(ctx context.Context, t T) error {
	conn := boiltx.DB(ctx, br.DB)

	entry, err := br.Mapper.ToEntry(t)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	err = entry.Insert(ctx, conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	return nil
}

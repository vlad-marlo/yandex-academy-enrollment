package pgx

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx"
	"go.uber.org/zap"
)

var (
	_               production.Store = (*Store)(nil)
	ErrNilReference                  = errors.New("unexpectedly got nil reference in storage")
)

// Store is postgres storage.
type Store struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

// New return storage with provided client
func New(cli pgx.Client) (*Store, error) {
	if cli == nil {
		return nil, ErrNilReference
	}
	return &Store{
		log:  cli.L(),
		pool: cli.P(),
	}, nil
}

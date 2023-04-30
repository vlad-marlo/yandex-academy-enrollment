package client

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Client struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

//func New() (*Client, error) {
//
//}

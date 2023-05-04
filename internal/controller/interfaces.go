package controller

import "context"

type Config interface {
	BindAddr() string
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

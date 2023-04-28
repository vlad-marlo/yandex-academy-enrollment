package controller

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
)

type Config interface {
	BindAddr() string
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Service interface {
	GetCourierByID(ctx context.Context, id int) (*model.Courier, error)
	CreateCouriers(ctx context.Context, request *model.CouriersCreateRequest) (*model.CouriersCreateResponse, error)
}

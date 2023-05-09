package production

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"go.uber.org/zap"
)

//go:generate mockgen --source=service.go --destination=mocks/service.go --package=mocks

type Store interface {
	GetCourierByID(ctx context.Context, id int) (*model.CourierDTO, error)
	CreateCouriers(ctx context.Context, couriers []model.CreateCourierDTO) ([]model.CourierDTO, error)
	GetCouriers(ctx context.Context, limit int, offset int) ([]model.CourierDTO, error)
	GetOrderByID(ctx context.Context, id int) (*model.OrderDTO, error)
	GetOrders(ctx context.Context, limit int, offset int) ([]*model.OrderDTO, error)
	CreateOrders(ctx context.Context, orders []*model.OrderDTO) error
}

type Service struct {
	log     *zap.Logger
	storage Store
}

func New(log *zap.Logger, storage Store) (*Service, error) {
	if log == nil || storage == nil {
		return nil, ErrNilReference
	}
	s := &Service{
		log:     log,
		storage: storage,
	}
	return s, nil
}

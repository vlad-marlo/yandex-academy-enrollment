package production

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen --source=service.go --destination=mocks/service.go --package=mocks

// Store storing all scrap.
type Store interface {
	// Courier methods

	GetCourierByID(ctx context.Context, id int64) (*model.CourierDTO, error)
	CreateCouriers(ctx context.Context, couriers []model.CreateCourierDTO) ([]model.CourierDTO, error)
	GetCouriers(ctx context.Context, limit int, offset int) ([]model.CourierDTO, error)

	// Order methods

	GetOrderByID(ctx context.Context, id int64) (*model.OrderDTO, error)
	GetOrders(ctx context.Context, limit int, offset int) ([]*model.OrderDTO, error)
	CreateOrders(ctx context.Context, orders []*model.OrderDTO) error
	GetCompletedOrdersPriceByCourier(ctx context.Context, id int64, start time.Time, end time.Time) ([]int32, error)
	CompleteOrders(ctx context.Context, info []model.CompleteOrder) error
	GetOrdersByIDs(ctx context.Context, ids []int64) ([]*model.OrderDTO, error)
}

var _ controller.Service = (*Service)(nil)

type Service struct {
	log     *zap.Logger
	storage Store
}

// New returns service with provided params.
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

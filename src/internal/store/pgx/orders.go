package pgx

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"time"
)

func (s *Store) GetOrderByID(ctx context.Context, id int64) (*model.OrderDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetOrders(ctx context.Context, limit int, offset int) ([]*model.OrderDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) CreateOrders(ctx context.Context, orders []*model.OrderDTO) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetCompletedOrdersPriceByCourier(_ context.Context, id int64, start time.Time, end time.Time) ([]int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) CompleteOrders(ctx context.Context, info []model.CompleteOrder) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetOrdersByIDs(ctx context.Context, ids []int64) ([]*model.OrderDTO, error) {
	//TODO implement me
	panic("implement me")
}

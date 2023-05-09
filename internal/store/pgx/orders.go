package pgx

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
)

func (s *Store) GetOrderByID(ctx context.Context, id int) (*model.OrderDTO, error) {
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

package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/multierr"
	"time"
)

func (s *Store) getDeliveryHoursOfOrder(ctx context.Context, id int64) (r []*datetime.TimeInterval, err error) {
	const query = `SELECT x.start_time, x.end_time, x.reversed
FROM orders_delivery_hours x
WHERE x.order_id = $1;`
	var rows pgx.Rows

	rows, err = s.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("pgxpool: doing query: %w", err)
	}
	defer rows.Close()
	var h datetime.TimeIntervalAlias

	for rows.Next() {
		if err = rows.Scan(h.Start, h.End, h.Reverse); err != nil {
			return nil, fmt.Errorf("error while scanning from rows: %w", err)
		}

		r = append(r, h.TimeInterval())
	}

	if err = rows.Err(); err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, pgx.ErrNoRows) || (errors.As(err, &pgErr) && pgerrcode.IsNoData(pgErr.Code)) {
			return r, nil
		}

		return nil, fmt.Errorf("error from rows.Err() => %w", err)
	}

	return r, nil
}

func (s *Store) GetOrderByID(ctx context.Context, id int64) (o *model.OrderDTO, err error) {
	const query = `SELECT x.weight, x.regions, x.cost, x.completed_time
FROM orders x
WHERE x.id = $1;`
	o = new(model.OrderDTO)

	if err = s.pool.QueryRow(ctx, query, id).Scan(&o.Weight, &o.Regions, &o.Cost, &o.CompletedTime); err != nil {
		return nil, fmt.Errorf("pgxpool: scan: %w", err)
	}

	if o.DeliveryHours, err = s.getDeliveryHoursOfOrder(ctx, id); err != nil {
		return nil, fmt.Errorf("get delivery hours of order: %w", err)
	}

	return o, nil
}

func (s *Store) getOrders(ctx context.Context, ids []int64) (res []*model.OrderDTO, err error) {
	res = make([]*model.OrderDTO, 0, len(ids))
	var courier *model.OrderDTO
	for _, id := range ids {
		courier, err = s.GetOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, courier)
	}
	return res, nil
}

func (s *Store) GetOrders(ctx context.Context, limit int, offset int) (res []*model.OrderDTO, err error) {
	const query = `SELECT x.id
FROM orders x
ORDER BY x.id
OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY;`
	var rows pgx.Rows

	ids := make([]int64, 0, limit)

	rows, err = s.pool.Query(ctx, query, offset, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNoContent
		}
		return nil, fmt.Errorf("unable to get users: %w", err)
	}

	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("unable to scan id: %w", err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("err from rows.Err(): %w", err)
	}

	res, err = s.getOrders(ctx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Store) createOrder(ctx context.Context, tx pgx.Tx, order *model.OrderDTO) (err error) {
	const query = `INSERT INTO orders(weight, regions, cost, completed)
VALUES ($1, $2, $3, FALSE)
RETURNING id;`
	multierr.AppendInto(&err, tx.QueryRow(ctx, query, order.Weight, order.Regions, order.Cost).Scan(&order.OrderID))
	return
}

func (s *Store) CreateOrders(ctx context.Context, orders []*model.OrderDTO) (err error) {
	var tx pgx.Tx
	tx, err = s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("check drivers: unable to begin tx: %w", err)
	}
	defer multierr.AppendInto(&err, tx.Rollback(ctx))
	for _, order := range orders {
		if multierr.AppendInto(&err, s.createOrder(ctx, tx, order)) {
			return err
		}
	}
	multierr.AppendInto(&err, tx.Commit(ctx))
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *Store) GetCompletedOrdersPriceByCourier(ctx context.Context, id int64, start time.Time, end time.Time) (sum, count int32, err error) {
	const query = `SELECT SUM(x.cost), COUNT(x.cost)
FROM orders x
WHERE x.completed
  AND x.courier = $1
  AND x.completed_time >= $2::TIMESTAMP
  AND x.completed_time <= $3::TIMESTAMP
ORDER BY id;`

	if err = s.pool.QueryRow(ctx, query, id, start, end).Scan(&sum, &count); err != nil {
		return 0, 0, err
	}
	return
}

func (s *Store) CompleteOrders(ctx context.Context, info []model.CompleteOrder) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetOrdersByIDs(ctx context.Context, ids []int64) (res []*model.OrderDTO, err error) {
	const query = `SELECT x.weight, x.regions, x.cost, x.completed_time
FROM orders x
WHERE x.id = $1;`
	res = make([]*model.OrderDTO, 0, len(ids))
	var order model.OrderDTO
	for _, id := range ids {
		if err = s.pool.QueryRow(ctx, query, id).Scan(&order.Weight, &order.Regions, &order.Cost, &order.CompletedTime); err != nil {
			return nil, err
		}

		if order.DeliveryHours, err = s.getDeliveryHoursOfOrder(ctx, id); err != nil {
			return nil, err
		}
	}
	return res, nil
}

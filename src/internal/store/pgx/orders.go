package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"time"
)

func (s *Store) getDeliveryHoursOfOrder(ctx context.Context, id int64) (r []*datetime.TimeInterval, err error) {
	const query = `SELECT x.start_time, x.end_time, x.reversed
FROM orders_delivery_hours x
WHERE x.order_id = $1;`
	var rows pgx.Rows

	r = make([]*datetime.TimeInterval, 0, 8)

	rows, err = s.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("err while doing query: %w", err)
	}
	defer rows.Close()

	var h datetime.TimeIntervalAlias

	for rows.Next() {
		if err = rows.Scan(&h.Start, &h.End, &h.Reverse); err != nil {
			return nil, fmt.Errorf("error while scanning from rows: %w", err)
		}

		r = append(r, h.TimeInterval())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error from rows.Err() => %w", err)
	}

	return r, nil
}

func (s *Store) GetOrderByID(ctx context.Context, id int64) (o *model.OrderDTO, err error) {
	const query = `SELECT x.weight, x.regions, x.cost, coalesce(x.completed_time, '1000-01-01'::timestamp)
FROM orders x
WHERE x.id = $1;`
	o = &model.OrderDTO{
		OrderID: id,
	}
	var t time.Time
	if err = s.pool.QueryRow(ctx, query, id).Scan(&o.Weight, &o.Regions, &o.Cost, &t); err != nil {
		return nil, fmt.Errorf("pgxpool: scan: %w", err)
	}
	o.CompletedTime = datetime.Time(t)
	if year, month, day := t.Date(); year == 1000 && month == time.January && day == 1 {
		o.CompletedTime = datetime.Time{}
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

func (s *Store) addDeliveryHoursToOrder(ctx context.Context, tx pgx.Tx, dto *model.OrderDTO) (err error) {
	for _, wh := range dto.DeliveryHours {
		if _, err = tx.Exec(
			ctx,
			`INSERT INTO orders_delivery_hours(order_id, start_time, end_time, reversed) VALUES ($1, $2, $3, $4);`,
			dto.OrderID,
			int(wh.Start()),
			int(wh.End()),
			wh.Start() > wh.End(),
		); err != nil {
			return fmt.Errorf("err while adding working hours to courier: %w", err)
		}
	}
	return nil
}

func (s *Store) createOrder(ctx context.Context, tx pgx.Tx, order *model.OrderDTO) (err error) {
	const query = `INSERT INTO orders(weight, regions, cost, completed)
VALUES ($1, $2, $3, FALSE)
RETURNING id;`

	if err = tx.QueryRow(ctx, query, order.Weight, order.Regions, order.Cost).Scan(&order.OrderID); err != nil {
		return fmt.Errorf("err while creating order: %w", err)
	}
	return s.addDeliveryHoursToOrder(ctx, tx, order)
}

func (s *Store) CreateOrders(ctx context.Context, orders []*model.OrderDTO) (err error) {
	var tx pgx.Tx
	tx, err = s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("check drivers: unable to begin tx: %w", err)
	}

	defer func() {
		s.log.Error("tx rollback", zap.NamedError("tx_error", tx.Rollback(ctx)))
	}()

	for _, order := range orders {
		if multierr.AppendInto(&err, s.createOrder(ctx, tx, order)) {
			return err
		}
	}

	if multierr.AppendInto(&err, tx.Commit(ctx)) {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *Store) GetCompletedOrdersPriceByCourier(ctx context.Context, id int64, start time.Time, end time.Time) (sum, count int32, err error) {
	const query = `SELECT COALESCE(SUM(x.cost), 0), COALESCE(COUNT(x.cost), 0)
FROM orders x
WHERE x.completed
  AND x.courier = $1
  AND x.completed_time >= $2::TIMESTAMP
  AND x.completed_time <= $3::TIMESTAMP;`

	if err = s.pool.QueryRow(ctx, query, id, start, end).Scan(&sum, &count); err != nil {
		return 0, 0, err
	}
	return
}

func (s *Store) completeOrder(ctx context.Context, tx pgx.Tx, order *model.CompleteOrder) error {
	var orderAssigned, orderCompleted bool
	const query = `SELECT EXISTS(
    SELECT * FROM orders WHERE courier = $1 AND id = $2
), EXISTS(
    SELECT * FROM orders x WHERE x.id = $2 AND completed = false
);`
	if err := tx.QueryRow(ctx, query, order.CourierID, order.OrderID).Scan(&orderAssigned, &orderCompleted); err != nil {
		return err
	}
	if !orderAssigned {
		return store.ErrDoesNotExists
	}
	if orderCompleted {
		return tx.QueryRow(ctx, `select completed_time from orders where id = $1;`, order.OrderID).Scan(&order.CompleteTime)
	}
	const updateQuery = `update orders set completed_time = $1, completed = TRUE where id = $2;`
	_, err := tx.Exec(ctx, updateQuery, order.CompleteTime, order.OrderID)
	return err
}

func (s *Store) CompleteOrders(ctx context.Context, info []model.CompleteOrder) (err error) {
	var tx pgx.Tx
	tx, err = s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		s.log.Error("rollback", zap.NamedError("tx_error", tx.Rollback(ctx)))
	}()

	for _, o := range info {
		if err = s.completeOrder(ctx, tx, &o); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) GetOrdersByIDs(ctx context.Context, ids []int64) (res []*model.OrderDTO, err error) {
	const query = `SELECT x.weight, x.regions, x.cost, coalesce(x.completed_time, '1000-01-01'::timestamp), x.completed
FROM orders x
WHERE x.id = $1;`
	res = make([]*model.OrderDTO, 0, len(ids))
	var (
		t     time.Time
		ok    bool
		order *model.OrderDTO
	)
	for _, id := range ids {
		order = new(model.OrderDTO)
		if err = s.pool.QueryRow(ctx, query, id).Scan(&order.Weight, &order.Regions, &order.Cost, &t, &ok); err != nil {
			return nil, err
		}
		if ok {
			order.CompletedTime = datetime.Time(t)
		}
		order.OrderID = id

		if order.DeliveryHours, err = s.getDeliveryHoursOfOrder(ctx, id); err != nil {
			return nil, err
		}
		res = append(res, order)
	}
	return res, nil
}

package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/multierr"
)

// getCourierRegions return slice of regions of courier with provided id.
func (s *Store) getCourierRegions(ctx context.Context, id int64) (r []int32, err error) {
	var rows pgx.Rows

	rows, err = s.pool.Query(ctx, `SELECT x.region FROM courier_region x WHERE x.courier_id = $1;`, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []int32{}, nil
		}

		return nil, fmt.Errorf("err while doing query: %w", err)
	}

	var i int32
	for rows.Next() {

		if err = rows.Scan(&i); err != nil {
			return nil, fmt.Errorf("error while scanning from rows: %w", err)
		}

		r = append(r, i)
	}

	if err = rows.Err(); err != nil && !(errors.Is(err, pgx.ErrNoRows)) {
		return nil, fmt.Errorf("error from rows.Err() => %w", err)
	}

	return r, nil
}

func (s *Store) getCourierWorkingHours(ctx context.Context, id int64) (r []*datetime.TimeInterval, err error) {
	var rows pgx.Rows

	r = []*datetime.TimeInterval{}

	rows, err = s.pool.Query(ctx, `SELECT x.start_time, x.end_time, x.reversed FROM courier_working_hour x WHERE x.courier_id = $1;`, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return r, nil
		}
		return nil, fmt.Errorf("err while doing ")
	}

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

func (s *Store) GetCourierByID(ctx context.Context, id int64) (courier *model.CourierDTO, err error) {
	courier = &model.CourierDTO{
		CourierID: id,
	}

	if err = s.pool.QueryRow(
		ctx,
		`SELECT x.courier_type FROM couriers x WHERE x.id=$1;`,
		id,
	).Scan(&courier.CourierType); err != nil {
		return nil, fmt.Errorf("unknown err while scanning: %w", err)
	}

	courier.Regions, err = s.getCourierRegions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unknown err while getting regions: %w", err)
	}
	courier.WorkingHours, err = s.getCourierWorkingHours(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unknown err while getting working hours: %w", err)
	}

	return courier, nil
}

func (s *Store) addRegionsToCourier(ctx context.Context, tx pgx.Tx, courier model.CourierDTO) (err error) {
	for _, region := range courier.Regions {
		if _, err = tx.Exec(
			ctx,
			`INSERT INTO courier_region(region, courier_id) VALUES ($1, $2);`,
			region,
			courier.CourierID,
		); err != nil {
			return fmt.Errorf("err while adding courier regions")
		}
	}
	return nil
}

func (s *Store) addWorkingHoursToCourier(ctx context.Context, tx pgx.Tx, dto model.CourierDTO) (err error) {
	for _, wh := range dto.WorkingHours {
		if _, err = tx.Exec(
			ctx,
			`INSERT INTO courier_working_hour(courier_id, start_time, end_time, reversed) VALUES ($1, $2, $3, $4);`,
			dto.CourierID,
			wh.Start(),
			wh.End(),
			wh.Start() > wh.End(),
		); err != nil {
			return fmt.Errorf("err while adding courier regions: %w", err)
		}
	}
	return nil
}

func (s *Store) createCourier(ctx context.Context, tx pgx.Tx, courier model.CreateCourierDTO) (model.CourierDTO, error) {
	res := model.CourierDTO{
		CourierID:    0,
		CourierType:  courier.CourierType,
		Regions:      courier.Regions,
		WorkingHours: courier.WorkingHours,
	}
	if err := tx.QueryRow(
		ctx,
		`insert into couriers(courier_type) values ($1) returning id;`,
		courier.CourierType,
	).Scan(&res.CourierID); err != nil {
		return model.CourierDTO{}, err
	}
	if err := s.addWorkingHoursToCourier(ctx, tx, res); err != nil {
		return model.CourierDTO{}, fmt.Errorf("error while adding WH to courier: %w", err)
	}
	if err := s.addRegionsToCourier(ctx, tx, res); err != nil {
		return model.CourierDTO{}, fmt.Errorf("error while adding regions to courier: %w", err)
	}

	return res, nil
}

func (s *Store) CreateCouriers(ctx context.Context, couriers []model.CreateCourierDTO) (r []model.CourierDTO, err error) {
	var tx pgx.Tx

	tx, err = s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: check drivers: %w", err)
	}
	defer multierr.AppendInto(&err, tx.Rollback(ctx))
	r = make([]model.CourierDTO, 0, len(couriers))

	var c model.CourierDTO
	for _, courier := range couriers {
		c, err = s.createCourier(ctx, tx, courier)
		if err != nil {
			return nil, fmt.Errorf("error while creating courier: %w", err)
		}
		r = append(r, c)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error while committing: update drivers: %w", err)
	}

	return r, nil
}

func (s *Store) getCouriers(ctx context.Context, ids []int64) (res []model.CourierDTO, err error) {
	res = make([]model.CourierDTO, 0, len(ids))
	var courier *model.CourierDTO
	for _, id := range ids {
		courier, err = s.GetCourierByID(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, *courier)
	}
	return res, nil
}

func (s *Store) GetCouriers(ctx context.Context, limit int, offset int) (res []model.CourierDTO, err error) {
	const query = `SELECT x.id FROM couriers x ORDER BY x.id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY;`
	var rows pgx.Rows

	ids := make([]int64, 0, limit)

	rows, err = s.pool.Query(ctx, query, offset, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.CourierDTO{}, nil
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

	res, err = s.getCouriers(ctx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

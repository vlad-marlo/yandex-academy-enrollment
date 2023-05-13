package migrator

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/retryer"
	"time"
)

const (
	migrationRetryAttempts = 2
	migrationsRetryDelay   = time.Second
)

var (
	migrations = []string{
		`CREATE TABLE IF NOT EXISTS couriers
(
    id           BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    courier_type TEXT                         NOT NULL,
    CONSTRAINT couriers_courier_type_check
        CHECK ((courier_type = 'FOOT'::TEXT) OR (courier_type = 'AUTO'::TEXT) OR
               (courier_type = 'BIKE'::TEXT))
);`,
		`CREATE TABLE IF NOT EXISTS courier_region
(
    id         BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    region     BIGINT                       NOT NULL,
    courier_id BIGINT                       NOT NULL,
    CONSTRAINT courier_fk FOREIGN KEY (courier_id) REFERENCES couriers MATCH FULL
);`,
		`CREATE TABLE IF NOT EXISTS courier_working_hour
(
    id         BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    start_time INTEGER DEFAULT 0            NOT NULL,
    end_time   INTEGER DEFAULT 0            NOT NULL,
    reversed   BOOLEAN DEFAULT FALSE        NOT NULL,
    courier_id BIGINT                       NOT NULL,
    CONSTRAINT courier_fk FOREIGN KEY (courier_id) REFERENCES couriers
);`,
		`CREATE TABLE IF NOT EXISTS orders
(
    id             BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    weight         FLOAT8                       NOT NULL,
    regions        INT4                         NOT NULL,
    cost           INT4                         NOT NULL,
    completed_time TIMESTAMP                    NULL,
    courier        BIGINT                       NULL,
    completed      BOOLEAN                      NOT NULL DEFAULT FALSE,
    CONSTRAINT check_complete CHECK ( (completed AND completed_time IS NOT NULL) OR
                                      (NOT completed AND completed_time IS NULL)),
    CONSTRAINT courier_fk FOREIGN KEY (courier) REFERENCES couriers (id) MATCH FULL,
    CONSTRAINT check_numerics CHECK ( cost > 0::INT4 AND regions > 0::INT4 AND weight > 0::FLOAT8 )
);`,
		`CREATE TABLE IF NOT EXISTS orders_delivery_hours
(
    id         BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    start_time INTEGER DEFAULT 0            NOT NULL,
    end_time   INTEGER DEFAULT 0            NOT NULL,
    reversed   BOOLEAN DEFAULT FALSE        NOT NULL,
    order_id   BIGINT                       NOT NULL,
    CONSTRAINT order_fk FOREIGN KEY (order_id) REFERENCES orders (id)
);`,
	}
)

func Migrate(cli pgx.Client) (int, error) {
	i := 0
	for _, migration := range migrations {
		if err := retryer.TryWithAttempts(
			func() error {
				_, err := cli.P().Exec(context.Background(), migration)
				return err
			},
			migrationRetryAttempts,
			migrationsRetryDelay,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.InvalidColumnDefinition {
					continue
				}
			}
			return i, err
		}
		i++
	}
	return i, nil
}

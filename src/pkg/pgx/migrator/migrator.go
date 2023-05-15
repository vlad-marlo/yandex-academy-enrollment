package migrator

import (
	"context"
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
    start_time INT4    DEFAULT 0::INT4      NOT NULL,
    end_time   INT4    DEFAULT 0::INT4      NOT NULL,
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
    start_time INT4    DEFAULT 0::INT4      NOT NULL,
    end_time   INT4    DEFAULT 0::INT4      NOT NULL,
    reversed   BOOLEAN DEFAULT FALSE        NOT NULL,
    order_id   BIGINT                       NOT NULL,
    CONSTRAINT order_fk FOREIGN KEY (order_id) REFERENCES orders (id)
);`,
		`CREATE TABLE IF NOT EXISTS order_group
(
    id      BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    date    VARCHAR(10)                  NOT NULL,
    courier BIGINT                       NOT NULL,
    CONSTRAINT courier_fk FOREIGN KEY (courier) REFERENCES couriers MATCH FULL ON DELETE CASCADE,
    CONSTRAINT courier_date_unique UNIQUE (date, courier)
);`,
	}
	migrateDown = []string{
		`DROP TABLE IF EXISTS order_group;`,
		`DROP TABLE IF EXISTS orders_delivery_hours;`,
		`DROP TABLE IF EXISTS orders;`,
		`DROP TABLE IF EXISTS courier_working_hour;`,
		`DROP TABLE IF EXISTS courier_region;`,
		`DROP TABLE IF EXISTS couriers;`,
	}
	Migrations = len(migrations)
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
			return i, err
		}
		i++
	}
	return i, nil
}

func MigrateDown(cli pgx.Client) (int, error) {
	i := 0
	for _, migration := range migrateDown {
		if err := retryer.TryWithAttempts(
			func() error {
				_, err := cli.P().Exec(context.Background(), migration)
				return err
			},
			migrationRetryAttempts,
			migrationsRetryDelay,
		); err != nil {
			return i, err
		}
		i++
	}
	return i, nil
}

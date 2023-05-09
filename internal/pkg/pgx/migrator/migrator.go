package migrator

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/pgx"
)

var (
	migrations = []string{
		`CREATE TABLE IF NOT EXISTS couriers
(
    id           BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
    courier_type TEXT                         NOT NULL,
    CONSTRAINT couriers_courier_type_check
        check ((courier_type = 'FOOT'::text) OR (courier_type = 'AUTO'::text) OR
               (courier_type = 'BIKE'::text))
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
    end_time   INTEGER default 0            NOT NULL,
    reversed   BOOLEAN DEFAULT FALSE        NOT NULL,
    courier_id BIGINT                       NOT NULL,
    CONSTRAINT courier_fk FOREIGN KEY (courier_id) REFERENCES couriers
);`,
	}
)

func Migrate(cli pgx.Client) (int, error) {
	i := 0
	for _, migration := range migrations {
		if _, err := cli.P().Exec(context.Background(), migration); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.DuplicateDatabase {
					continue
				}
			}
			return i, err
		}
		i++
	}
	return i, nil
}

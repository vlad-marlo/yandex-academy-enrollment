package migrator

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/pgx"
	"go.uber.org/zap"
)

const (
	createVersionTableTemplate = `CREATE TABLE IF NOT EXISTS migrate_version(
    version int
);`
	getVersionTableTemplate       = `SELECT version FROM migrate_version;`
	versionExistsTemplate         = "SELECT EXISTS(SELECT * FROM migrate_version);"
	incrementVersionTableTemplate = `UPDATE migrate_version SET version = version + $1;`
)

var (
	migrations = []string{
		`CREATE TABLE IF NOT EXISTS time
(
    id     bigserial unique primary key not null,
    hour   int,
    minute int
);
CREATE TABLE IF NOT EXISTS time_interval
(
    id         bigserial unique primary key not null,
    start_time bigint references time (id)  not null,
    end_time   bigint references time (id)  not null,
    reverse    bool default false
);`,
		`CREATE TYPE courier_type AS ENUM (
    'FOOT',
    'BIKE',
    'AUTO'
    );
CREATE TABLE IF NOT EXISTS couriers
(
    id   bigserial unique primary key not null,
    type courier_type
);`,
	}
)

type migrator struct {
	pool    *pgxpool.Pool
	log     *zap.Logger
	ctx     context.Context
	version int
}

func Migrate(cli pgx.Client) error {
	m := &migrator{
		ctx:  context.Background(),
		log:  cli.L(),
		pool: cli.P(),
	}
	return m.migrate()
}

func (m *migrator) migrate() error {
	if err := m.createVersionTable(); err != nil {
		return err
	}
	return nil
}

func (m *migrator) createVersionTable() error {
	_, err := m.pool.Exec(
		m.ctx,
		createVersionTableTemplate,
	)
	if err != nil {
		return fmt.Errorf("err while creating table: %w", err)
	}
	var ok bool
	if err = m.pool.QueryRow(
		m.ctx,
		versionExistsTemplate,
	).Scan(&ok); err != nil {
		return err
	}
	return err
}

func (m *migrator) getDatabaseVersion() (ver int, err error) {
	if err = m.pool.QueryRow(m.ctx, getVersionTableTemplate).Scan(&ver); err != nil {
		return 0, err
	}
	return ver, nil
}

func (m *migrator) incrementDatabaseVersion() error {
	_, err := m.pool.Exec(m.ctx, incrementVersionTableTemplate)
	return err
}

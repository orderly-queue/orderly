package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/orderly-queue/orderly/database/migrations"
)

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(db *sql.DB) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	fs, err := iofs.New(migrations.Migrations, "files")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithInstance(
		"iofs",
		fs,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}
	return &Migrator{
		m: m,
	}, nil
}

func (m *Migrator) Up() error {
	err := m.m.Up()
	if err == nil || err.Error() == "no change" {
		return nil
	}
	return err
}

func (m *Migrator) Down() error {
	err := m.m.Down()
	if err == nil || err.Error() == "no change" {
		return nil
	}
	return err
}

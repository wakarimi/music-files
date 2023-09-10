package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func ConnectDb(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	log.Debug().Msg("Migrations running")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create driver")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create migration")
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("Failed to apply migration")
		return err
	}

	log.Debug().Msg("Migrations applied successfully")
	return nil
}

package storage

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/config"
	"strings"
)

func New(cfg config.DBConfig) (db *sqlx.DB, err error) {
	db, err = connectDB(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	err = runMigrations(db, cfg.MigrationPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute migration")
		return nil, err
	}

	return db, nil
}

func connectDB(cfg config.DBConfig) (db *sqlx.DB, err error) {
	connectionStringBuilder := strings.Builder{}
	connectionStringBuilder.WriteString(fmt.Sprintf("postgresql://%s", cfg.User))
	connectionStringBuilder.WriteString(fmt.Sprintf(":%s", cfg.Password))
	connectionStringBuilder.WriteString(fmt.Sprintf("@%s", cfg.Host))
	connectionStringBuilder.WriteString(fmt.Sprintf(":%d", cfg.Port))
	connectionStringBuilder.WriteString(fmt.Sprintf("/%s", cfg.DBName))
	connectionStringBuilder.WriteString("?sslmode=disable")
	connectionString := connectionStringBuilder.String()

	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping database")
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sqlx.DB, migrationsPath string) error {
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

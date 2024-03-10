package postgres

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"music-files/internal/config"
	"strings"
)

type Postgres struct {
	*sqlx.DB
}

func New(cfg config.DBConfig) (*Postgres, error) {
	db, err := connectDB(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	err = runMigrations(db, cfg.MigrationPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}

	return &Postgres{
		DB: db,
	}, nil
}

func connectDB(cfg config.DBConfig) (db *sqlx.DB, err error) {
	connectionStringBuilder := strings.Builder{}
	connectionStringBuilder.WriteString(fmt.Sprintf("postgresql://%s", cfg.Username))
	connectionStringBuilder.WriteString(fmt.Sprintf(":%s", cfg.Password))
	connectionStringBuilder.WriteString(fmt.Sprintf("@%s", cfg.Host))
	connectionStringBuilder.WriteString(fmt.Sprintf(":%d", cfg.Port))
	connectionStringBuilder.WriteString(fmt.Sprintf("/%s", cfg.Name))
	connectionStringBuilder.WriteString("?sslmode=disable")
	connectionString := connectionStringBuilder.String()

	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping connected database")
	}

	return db, nil
}

func runMigrations(db *sqlx.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create driver")
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create migration")
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to apply migration")
	}

	return nil
}

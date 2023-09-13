package context

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/config"
)

type AppContext struct {
	Config *config.Configuration
	Db     *sqlx.DB
}

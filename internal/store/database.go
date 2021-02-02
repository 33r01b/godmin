package store

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"godmin/config"
)

func NewDB(config *config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dbUrl(config))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database: %v")
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return db, nil
}

func dbUrl(config *config.Database) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SslMode,
	)
}

package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/store"
)

func Run(context context.Context, config *config.Config) error {
	log.SetLevel(config.LogLevel)

	db, err := store.NewDB(&config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

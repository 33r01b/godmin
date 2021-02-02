package server

import (
	"context"
	"godmin/config"
	"godmin/internal/store"
)

func Run(context context.Context, config *config.Config) error {
	db, err := store.NewDB(&config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server/api"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.TODO()
	return api.Run(ctx, config.NewConfig())
}

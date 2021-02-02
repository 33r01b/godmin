package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	config.Bootstrap()
}

func run() error {
	ctx := context.TODO()
	return server.Run(ctx, config.NewConfig())
}

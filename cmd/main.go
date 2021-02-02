package main

import (
	"context"
	"godmin/config"
	"godmin/internal/server"
	"log"
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

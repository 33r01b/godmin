package main

import (
	"godmin/config"
	"godmin/internal/server"
	"godmin/internal/server/api"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.NewConfig()
	connections, err := server.NewConnections(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer connections.Close()

	apiServer := api.NewApi(conf, api.NewServices(connections, conf))
	apiServer.Run()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	select {
	case x := <-sigc:
		log.Info("received a signal.", x.String())
	case err := <-apiServer.Notify():
		log.Error("received an error from the api server.", "err", err)
	}

	if err := apiServer.Shutdown(); err != nil {
		log.Error(err)
	}
}

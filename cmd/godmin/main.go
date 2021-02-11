package main

import (
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server"
	"godmin/internal/server/api"
	"os"
	"os/signal"
	"syscall"
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		log.Info("received a signal.", "signal", x.String())
	case err := <-apiServer.Notify():
		log.Error("received an error from the api server.", "err", err)
	}

	if err := apiServer.Shutdown(); err != nil {
		log.Error(err)
	}
}

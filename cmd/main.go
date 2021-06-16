package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/common/cmd"
	"orderservice/pkg/order/infrastructure/transport"
)

const appID = "order"

type config struct {
	cmd.WebConfig
	cmd.DatabaseConfig
}

func main() {
	var conf config
	if err := envconfig.Process(appID, &conf); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	srv := startServer(&conf)

	cmd.WaitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(conf *config) *http.Server {
	log.WithFields(log.Fields{"port": conf.ServerPort}).Info("starting the order server")
	db := cmd.CreateDBConnection(conf.DatabaseConfig)

	router := transport.Router(db)

	srv := &http.Server{Addr: fmt.Sprintf(":%s", conf.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}

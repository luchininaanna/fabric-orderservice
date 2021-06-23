package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/api/storeservice"
	"orderservice/pkg/common/cmd"
	transportUtils "orderservice/pkg/common/infrastructure/transport"
	"orderservice/pkg/order/infrastructure/transport"
)

const appID = "order"

type config struct {
	cmd.WebConfig
	cmd.DatabaseConfig

	StoreGRPCAddress string `envconfig:"store_grpc_address"`
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
	defer transportUtils.CloseService(db, "database connection")

	storeConn := transportUtils.DialGRPC(conf.StoreGRPCAddress)
	defer transportUtils.CloseService(storeConn, "scoring connection")

	router := transport.Router(db, storeservice.NewStoreServiceClient(storeConn))

	srv := &http.Server{Addr: fmt.Sprintf(":%s", conf.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}

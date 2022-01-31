package main

import (
	"github.com/allinbits/emeris-cns-server/cns/auth"
	"github.com/allinbits/emeris-cns-server/cns/chainwatch"
	"github.com/allinbits/emeris-cns-server/cns/config"
	"github.com/allinbits/emeris-cns-server/cns/database"
	"github.com/allinbits/emeris-cns-server/cns/rest"
	"github.com/allinbits/emeris-utils/k8s"
	"github.com/allinbits/emeris-utils/logging"
)

var Version = "not specified"

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.New(logging.LoggingConfig{
		LogPath: config.LogPath,
		Debug:   config.Debug,
	})

	logger.Infow("cns-server", "version", Version)

	di, err := database.New(config.DatabaseConnectionURL)
	if err != nil {
		logger.Fatal(err)
	}

	kube, err := k8s.NewInCluster()
	if err != nil {
		logger.Fatal(err)
	}

	rc, err := chainwatch.NewConnection(config.Redis)
	if err != nil {
		logger.Fatal(err)
	}

	err = auth.NewOAuthServer(config.Env)
	if err != nil {
		logger.Fatal(err)
	}

	ci := chainwatch.New(
		logger,
		kube,
		config.KubernetesNamespace,
		rc,
		di,
		config.RelayerDebug,
	)

	go ci.Run()

	restServer := rest.NewServer(
		logger,
		di,
		kube,
		rc,
		config,
	)

	if err := restServer.Serve(config.RESTAddress); err != nil {
		logger.Panicw("rest http server error", "error", err)
	}
}

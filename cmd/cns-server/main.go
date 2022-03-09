package main

import (
	"github.com/emerishq/emeris-cns-server/cns/auth"
	"github.com/emerishq/emeris-cns-server/cns/chainwatch"
	"github.com/emerishq/emeris-cns-server/cns/config"
	"github.com/emerishq/emeris-cns-server/cns/database"
	"github.com/emerishq/emeris-cns-server/cns/rest"
	"github.com/emerishq/emeris-utils/k8s"
	"github.com/emerishq/emeris-utils/logging"
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

	logger.Infow("cns-server", "version", Version, "redirectURL", config.RedirectURL)

	di, err := database.New(config.DatabaseConnectionURL)
	if err != nil {
		logger.Fatal(err)
	}

	k8sConfig, err := k8s.InClusterConfig()
	if err != nil {
		logger.Fatal(err)
	}

	kube, err := k8s.NewClient(k8sConfig)
	if err != nil {
		logger.Fatal(err)
	}

	rc, err := chainwatch.NewConnection(config.Redis)
	if err != nil {
		logger.Fatal(err)
	}

	a, err := auth.NewOAuthServer(config.Env, config.RedirectURL, config.OAuth2ClientID, config.OAuth2ClientSecret, []byte(config.Secret))
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
		a,
	)

	if err := restServer.Serve(config.RESTAddress); err != nil {
		logger.Panicw("rest http server error", "error", err)
	}
}

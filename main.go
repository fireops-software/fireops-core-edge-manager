package main

import (
	"strings"

	"github.com/fireops-software/fireops-core-edge-manager/api"
	"github.com/fireops-software/fireops-core-edge-manager/dal"
	"github.com/fireops-software/fireops-core-edge-manager/logic"
	"github.com/uoul/go-common/config"
	"github.com/uoul/go-common/log"
)

const (
	VERSION = "{VERSION}"
)

func main() {
	// Create ConfigProvider
	cp := config.NewEnvVarProvider()
	// Create Logger
	logger := log.NewConsoleLogger(
		log.StringToLogLevel(
			cp.StringOrDefault("LOG_LEVEL", "INFO"),
			log.INFO,
		),
	)
	// Create FireOpsClient
	fireOpsClient := dal.NewFireOpsApi(
		cp.StringOrDefault("FIREOPS_BASE_URL", ""),
	)
	// Create Logic
	appLogic := logic.NewLogic(
		logger,
		fireOpsClient,
	)
	// Extract Api keys
	apiKeys := strings.Split(
		cp.StringOrDefault("API_KEYS", ""),
		",",
	)
	// Create Api
	restApi := api.NewApi(
		logger,
		appLogic,
		api.WithApiKeys(apiKeys),
	)
	// Run Api
	listen := cp.StringOrDefault("API_INTERFACE", ":80")
	logger.Infof("Api running on %s", listen)
	restApi.Run(listen)
}

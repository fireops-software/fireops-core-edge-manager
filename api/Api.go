package api

import (
	"fmt"
	"net/http"

	"github.com/fireops-software/fireops-core-edge-manager/logic"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/servemux"
)

const (
	API_BASE_URI = "/api/v1"
)

type Api struct {
	logger  log.ILogger
	logic   logic.ILogic
	apiKeys []string
}

func (a *Api) Run(listen string) error {
	// Create serve mux
	mux := http.NewServeMux()
	// Register Handlers
	a.registerDeviceEnpoints(mux)
	// Run HTTP Server
	return http.ListenAndServe(listen, mux)
}

func (a *Api) registerDeviceEnpoints(mux *http.ServeMux) {
	// Handle websocket request
	mux.HandleFunc(fmt.Sprintf("%s/ws", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		a.handleWebsocketRequest(),
		useErrorTranslation[any](a.logger),
	))
	// GET connected devices
	mux.HandleFunc(fmt.Sprintf("GET %s/devices", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.getConnectedDevices())...,
	))
	// GET current agent's version
	mux.HandleFunc(fmt.Sprintf("GET %s/devices/{deviceId}/version", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.getDeviceVersion())...,
	))
	// GET containers
	mux.HandleFunc(fmt.Sprintf("GET %s/devices/{deviceId}/containers", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.getContainers())...,
	))
	// GET container logs
	mux.HandleFunc(fmt.Sprintf("GET %s/devices/{deviceId}/containers/{containerId}/logs", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.getContainerLogs())...,
	))
	// PUT containers
	mux.HandleFunc(fmt.Sprintf("PUT %s/devices/{deviceId}/containers", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.updateContainers())...,
	))
	// DELETE containers
	mux.HandleFunc(fmt.Sprintf("DELETE %s/devices/{deviceId}/containers", API_BASE_URI), servemux.Handle(
		servemux.NewHandlerConfig(),
		defaultChain(a, a.deleteContainers())...,
	))
}

func defaultChain[T any](api *Api, handler servemux.HandlerFunc[T]) []servemux.HandlerFunc[T] {
	return []servemux.HandlerFunc[T]{
		useContentTypeApplicationJson[T](),
		useApiKey[T](api.apiKeys),
		handler,
		useErrorTranslation[T](api.logger),
	}
}

func WithApiKeys(keys []string) func(*Api) {
	return func(a *Api) {
		a.apiKeys = keys
	}
}

func NewApi(logger log.ILogger, logic logic.ILogic, opts ...func(*Api)) *Api {
	a := &Api{
		logger:  logger,
		apiKeys: []string{},
		logic:   logic,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

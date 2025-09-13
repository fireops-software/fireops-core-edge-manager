package logic

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/fireops-software/fireops-core-edge-manager/api/ws"
	"github.com/fireops-software/fireops-core-edge-manager/dal"
	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/gorilla/websocket"
	"github.com/uoul/go-common/log"

	appError "github.com/fireops-software/fireops-core-edge-manager/error"
)

type Logic struct {
	logger     log.ILogger
	fireOpsApi dal.IFireOpsApi
	mux        sync.Mutex
	wsTimeout  time.Duration
	devices    map[string]ws.WsRequestClient
}

// RegisterDevice implements ILogic.
func (l *Logic) RegisterDevice(ctx context.Context, deviceId string, conn *websocket.Conn) (context.Context, error) {
	// Lock to prevent inconsistency
	l.mux.Lock()
	defer l.mux.Unlock()
	// Check if device already registered
	if _, exists := l.devices[deviceId]; exists {
		return nil, appError.NewErrConflict("device already exists")
	}
	// Register device
	wsClient := ws.NewWsRequestClient(
		ctx,
		conn,
	)
	l.devices[deviceId] = *wsClient
	return wsClient.Context(), nil
}

// UnRegisterDevice implements ILogic.
func (l *Logic) UnRegisterDevice(ctx context.Context, deviceId string) error {
	// Lock to prevent inconsistency
	l.mux.Lock()
	defer l.mux.Unlock()
	// Check if device is registered
	device, exists := l.devices[deviceId]
	if !exists {
		return appError.NewErrNotFound("device does not exist")
	}
	// Close channels
	device.Close()
	delete(l.devices, deviceId)
	return nil
}

// GetDeviceIdentityFromToken implements ILogic.
func (l *Logic) GetDeviceIdentityFromToken(ctx context.Context, token string) (domain.DeviceIdentity, error) {
	identity := <-l.fireOpsApi.GetDeviceIdentityForToken(ctx, token)
	return identity.Result, identity.Error
}

// DeleteAllContainers implements ILogic.
func (l *Logic) DeleteAllContainers(ctx context.Context, deviceId string) error {
	// Send destroy message
	panic("unimplemented")
}

// DeployContainers implements ILogic.
func (l *Logic) DeployContainers(ctx context.Context, deviceId string, containerConfig []domain.ServiceDefinition) ([]container.Summary, error) {
	panic("unimplemented")
}

// GetAgentVersion implements ILogic.
func (l *Logic) GetDeviceVersion(ctx context.Context, deviceId string) (domain.DeviceVersion, error) {
	panic("unimplemented")
}

// GetContainerLogs implements ILogic.
func (l *Logic) GetContainerLogs(ctx context.Context, deviceId string, containerId string) ([]domain.ContainerLogEntry, error) {
	panic("unimplemented")
}

// GetContainers implements ILogic.
func (l *Logic) GetContainers(ctx context.Context, deviceId string) ([]container.Summary, error) {
	panic("unimplemented")
}

func NewLogic(logger log.ILogger, opts ...func(*Logic)) ILogic {
	l := &Logic{
		logger:     logger,
		fireOpsApi: dal.NewFireOpsApi(),
		mux:        sync.Mutex{},

		wsTimeout: 300 * time.Second,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

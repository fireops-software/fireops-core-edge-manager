package logic

import (
	"context"
	"encoding/json"
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
	devices    map[string]struct {
		toWebsocket   chan any
		fromWebsocket chan ws.Response[json.RawMessage]
		conn          *websocket.Conn
	}
}

// ListenWebsocketRequest implements ILogic.
func (l *Logic) HandleWebsocketRequest(ctx context.Context, deviceId string) error {
	// Check if device is registered
	device, exists := l.devices[deviceId]
	if !exists {
		return appError.NewErrNotFound("device does not exist")
	}
	// Listen for Request
	req, ok := <-device.toWebsocket
	if !ok {
		return appError.NewErrChannelClosed("request channel has been closed")
	}
	// Send reqeust to websocket
	if err := device.conn.WriteJSON(req); err != nil {
		return appError.NewErrNetwork("failed to write data to websocket (deviceId: %s) - %v", deviceId, err)
	}
	// Wait for response on Websocket
	device.conn.SetReadDeadline(time.Now().Add(l.wsTimeout))
	resp := ws.Response[json.RawMessage]{}
	if err := device.conn.ReadJSON(resp); err != nil {
		return appError.NewErrNetwork("failed to read response (deviceId: %s) - %v", deviceId, err)
	}
	// Send response
	device.fromWebsocket <- resp
	return nil
}

// RegisterDevice implements ILogic.
func (l *Logic) RegisterDevice(ctx context.Context, deviceId string, conn *websocket.Conn) error {
	// Lock to prevent inconsistency
	l.mux.Lock()
	defer l.mux.Unlock()
	// Check if device already registered
	if _, exists := l.devices[deviceId]; exists {
		return appError.NewErrConflict("device already exists")
	}
	// Register device
	l.devices[deviceId] = struct {
		toWebsocket   chan any
		fromWebsocket chan ws.Response[json.RawMessage]
		conn          *websocket.Conn
	}{
		toWebsocket:   make(chan any),
		fromWebsocket: make(chan ws.Response[json.RawMessage]),
		conn:          conn,
	}
	return nil
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
	close(device.fromWebsocket)
	close(device.toWebsocket)
	return nil
}

// GetDeviceIdentityFromToken implements ILogic.
func (l *Logic) GetDeviceIdentityFromToken(ctx context.Context, token string) (domain.DeviceIdentity, error) {
	identity := <-l.fireOpsApi.GetDeviceIdentityForToken(ctx, token)
	return identity.Result, identity.Error
}

// DeleteAllContainers implements ILogic.
func (l *Logic) DeleteAllContainers(ctx context.Context, deviceId string) error {
	// Get device
	device, exists := l.devices[deviceId]
	if !exists {
		return appError.NewErrNotFound("No device with id %s registered", deviceId)
	}
	// Send destroy message
	req := 
	device.toWebsocket <- ws.Request[ws.DestroyContainersRequest]{

	}
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
		devices: map[string]struct {
			toWebsocket   chan any
			fromWebsocket chan ws.Response[json.RawMessage]
			conn          *websocket.Conn
		}{},

		wsTimeout: 300 * time.Second,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

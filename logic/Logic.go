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
	"github.com/google/uuid"
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
	l.logger.Debugf("RegisterDevice(%s)...", deviceId)
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
	l.logger.Debugf("UnRegisterDevice(%s)...", deviceId)
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
	l.logger.Debugf("GetDeviceIdentityFromToken(%s)...", token)
	identity := <-l.fireOpsApi.GetDeviceIdentityForToken(ctx, token)
	return identity.Result, identity.Error
}

// DeleteAllContainers implements ILogic.
func (l *Logic) DeleteAllContainers(ctx context.Context, deviceId string) error {
	l.logger.Debugf("DeleteAllContainers(%s)...", deviceId)
	_, err := doWsRequest[any](
		l.devices,
		deviceId,
		&ws.Request[ws.DestroyContainersRequest]{
			MsgId:   uuid.NewString(),
			MsgType: ws.TYPE_DESTROY,
			Body:    ws.DestroyContainersRequest{},
		},
	)
	return err
}

// DeployContainers implements ILogic.
func (l *Logic) DeployContainers(ctx context.Context, deviceId string, containerConfig []domain.ServiceDefinition) ([]container.Summary, error) {
	l.logger.Debugf("DeployContainers(%s): %v", deviceId, containerConfig)
	return doWsRequest[[]container.Summary](
		l.devices,
		deviceId,
		&ws.Request[ws.InstallContainersRequest]{
			MsgId:   uuid.NewString(),
			MsgType: ws.TYPE_INSTALL,
			Body:    containerConfig,
		},
	)
}

// GetAgentVersion implements ILogic.
func (l *Logic) GetDeviceVersion(ctx context.Context, deviceId string) (domain.DeviceVersion, error) {
	l.logger.Debugf("GetDeviceVersion(%s)...", deviceId)
	return doWsRequest[domain.DeviceVersion](
		l.devices,
		deviceId,
		&ws.Request[ws.GetAgentVersionRequest]{
			MsgId:   uuid.NewString(),
			MsgType: ws.TYPE_GET_AGENT_VERSION,
			Body:    ws.GetAgentVersionRequest{},
		},
	)
}

// GetContainerLogs implements ILogic.
func (l *Logic) GetContainerLogs(ctx context.Context, deviceId string, containerId string, len uint) ([]domain.ContainerLogEntry, error) {
	l.logger.Debugf("GetContainerLogs(%s)...", deviceId)
	return doWsRequest[[]domain.ContainerLogEntry](
		l.devices,
		deviceId,
		&ws.Request[ws.GetContainerLogsRequest]{
			MsgId:   uuid.NewString(),
			MsgType: ws.TYPE_GET_CONTAINER_LOGS,
			Body: ws.GetContainerLogsRequest{
				ContainerId: containerId,
				Len:         len,
			},
		},
	)
}

// GetContainers implements ILogic.
func (l *Logic) GetContainers(ctx context.Context, deviceId string) ([]container.Summary, error) {
	l.logger.Debugf("GetContainers(%s)...", deviceId)
	return doWsRequest[[]container.Summary](
		l.devices,
		deviceId,
		&ws.Request[ws.GetContainersRequest]{
			MsgId:   uuid.NewString(),
			MsgType: ws.TYPE_GET_CONTAINERS,
			Body:    ws.GetContainersRequest{},
		},
	)
}

func doWsRequest[T any](devices map[string]ws.WsRequestClient, deviceId string, req ws.IRequest) (T, error) {
	// Check if device is registered
	device, exists := devices[deviceId]
	if !exists {
		return *new(T), appError.NewErrNotFound("device with id %s not found", deviceId)
	}
	// Send deployment request
	resp := <-device.Send(req)
	if resp.Error != nil {
		return *new(T), resp.Error
	}
	if resp.Result.GetError() != nil {
		return *new(T), resp.Result.GetError()
	}
	// Convert Response
	respBody := *new(T)
	if err := json.Unmarshal(resp.Result.GetBody(), &respBody); err != nil {
		return *new(T), appError.NewErrDataParsing("failed to parse data - %v", err)
	}
	return respBody, nil
}

func NewLogic(logger log.ILogger, fireopsApi dal.IFireOpsApi, opts ...func(*Logic)) ILogic {
	l := &Logic{
		logger:     logger,
		fireOpsApi: fireopsApi,
		mux:        sync.Mutex{},
		devices:    map[string]ws.WsRequestClient{},

		wsTimeout: 300 * time.Second,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

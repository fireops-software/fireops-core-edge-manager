package logic

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/gorilla/websocket"
)

type ILogic interface {
	// Websocket Api
	GetDeviceIdentityFromToken(ctx context.Context, token string) (domain.DeviceIdentity, error)
	RegisterDevice(ctx context.Context, deviceId string, conn *websocket.Conn) (context.Context, error)
	UnRegisterDevice(ctx context.Context, deviceId string) error
	// REST Api
	GetDeviceVersion(ctx context.Context, deviceId string) (domain.DeviceVersion, error)
	GetContainers(ctx context.Context, deviceId string) ([]container.Summary, error)
	DeployContainers(ctx context.Context, deviceId string, containerConfig []domain.ServiceDefinition) ([]container.Summary, error)
	GetContainerLogs(ctx context.Context, deviceId string, containerId string, len uint) ([]domain.ContainerLogEntry, error)
	DeleteAllContainers(ctx context.Context, deviceId string) error
}

package ws

import "github.com/fireops-software/fireops-core-edge-manager/domain"

type GetContainerLogsRequest struct {
	ContainerId string
	Len         uint
}

type GetContainerLogsResponse []domain.ContainerLogEntry

package ws

import (
	"github.com/docker/docker/api/types/container"
	"github.com/fireops-software/fireops-core-edge-manager/domain"
)

type InstallContainersRequest []domain.ServiceDefinition

type InstallContainersResponse []container.Summary

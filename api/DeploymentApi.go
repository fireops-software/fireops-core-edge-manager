package api

import (
	"net/http"
	"strconv"

	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/uoul/go-common/servemux"
)

// Get connected devices
func (a *Api) getConnectedDevices() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		devices, err := a.logic.GetConnectedDevices(ctx.Context())
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.SetResponseBody(devices)
	}
}

// GET current agent's version
func (a *Api) getDeviceVersion() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		version, err := a.logic.GetDeviceVersion(
			ctx.Context(),
			ctx.GetUrlParam("deviceId"),
		)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.SetResponseBody(version)
	}
}

// GET containers
func (a *Api) getContainers() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		containers, err := a.logic.GetContainers(
			ctx.Context(),
			ctx.GetUrlParam("deviceId"),
		)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.SetResponseBody(containers)
	}
}

// GET container logs
func (a *Api) getContainerLogs() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		len, err := strconv.ParseUint(ctx.GetQueryParam("len"), 10, 32)
		if err != nil {
			len = 100
		}
		logs, err := a.logic.GetContainerLogs(
			ctx.Context(),
			ctx.GetUrlParam("deviceId"),
			ctx.GetUrlParam("containerId"),
			uint(len),
		)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.SetResponseBody(logs)
	}
}

// PUT containers
func (a *Api) updateContainers() servemux.HandlerFunc[[]domain.ServiceDefinition] {
	return func(ctx *servemux.HttpCtx[[]domain.ServiceDefinition]) {
		containerConfig, err := ctx.GetBody()
		if err != nil {
			ctx.Error(err)
			return
		}
		containers, err := a.logic.DeployContainers(
			ctx.Context(),
			ctx.GetUrlParam("deviceId"),
			containerConfig,
		)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetResponseBody(containers)
	}
}

// DELETE containers
func (a *Api) deleteContainers() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		err := a.logic.DeleteAllContainers(
			ctx.Context(),
			ctx.GetUrlParam("deviceId"),
		)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
}

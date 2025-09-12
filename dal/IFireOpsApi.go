package dal

import (
	"context"

	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/uoul/go-common/async"
)

type IFireOpsApi interface {
	GetDeviceIdentityForToken(ctx context.Context, token string) chan async.ActionResult[domain.DeviceIdentity]
}

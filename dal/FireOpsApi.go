package dal

import (
	"context"

	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/uoul/go-common/async"
)

type FireOpsApi struct{}

// GetAgentIdentityForToken implements IFireOpsApi.
func (f *FireOpsApi) GetDeviceIdentityForToken(ctx context.Context, token string) chan async.ActionResult[domain.DeviceIdentity] {
	r := make(chan async.ActionResult[domain.DeviceIdentity])
	go func() {
		// TODO: Implement API-Call for token validation
		r <- async.ActionResult[domain.DeviceIdentity]{
			Result: domain.DeviceIdentity{
				Id: "1",
			},
		}
	}()
	return r
}

func NewFireOpsApi() IFireOpsApi {
	return &FireOpsApi{}
}

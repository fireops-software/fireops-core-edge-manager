package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fireops-software/fireops-core-edge-manager/domain"
	"github.com/uoul/go-common/async"

	appError "github.com/fireops-software/fireops-core-edge-manager/error"
)

type FireOpsApi struct {
	timeout time.Duration
	baseUrl string
}

// GetAgentIdentityForToken implements IFireOpsApi.
func (f *FireOpsApi) GetDeviceIdentityForToken(ctx context.Context, token string) chan async.ActionResult[domain.DeviceIdentity] {
	r := make(chan async.ActionResult[domain.DeviceIdentity])
	go func() {
		// Create Request context
		ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
		defer cancel()
		// Create http request
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/edge/user/get-token-id", f.baseUrl), nil)
		if err != nil {
			r <- async.ActionResult[domain.DeviceIdentity]{
				Error: appError.NewErrNetwork("failed to create http request - %v", err),
			}
			return
		}

		// Add Headers
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		// Do request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			r <- async.ActionResult[domain.DeviceIdentity]{
				Error: appError.NewErrNetwork("failed to perform http request - %v", err),
			}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			r <- async.ActionResult[domain.DeviceIdentity]{
				Error: appError.NewErrNetwork("failed to validate token [StatusCode: %d]", resp.StatusCode),
			}
			return
		}
		// Parse Response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			r <- async.ActionResult[domain.DeviceIdentity]{
				Error: appError.NewErrDataParsing("failed read response body - %v", err),
			}
			return
		}
		identity := domain.DeviceIdentity{}
		if err := json.Unmarshal(body, &identity); err != nil {
			r <- async.ActionResult[domain.DeviceIdentity]{
				Error: appError.NewErrDataParsing("failed read response body - %v", err),
			}
			return
		}
		r <- async.ActionResult[domain.DeviceIdentity]{
			Result: identity,
			Error:  nil,
		}
	}()
	return r
}

func WithFireOpsApiTimeout(timeout time.Duration) func(*FireOpsApi) {
	return func(foa *FireOpsApi) {
		foa.timeout = timeout
	}
}

func NewFireOpsApi(baseUrl string, opts ...func(*FireOpsApi)) IFireOpsApi {
	f := &FireOpsApi{
		timeout: 10 * time.Second,
		baseUrl: baseUrl,
	}
	for _, o := range opts {
		o(f)
	}
	return f
}

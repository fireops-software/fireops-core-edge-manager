package api

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/uoul/go-common/servemux"

	appError "github.com/fireops-software/fireops-core-edge-manager/error"
)

// ----------------------------------------------------------------------
// Handler functions
// ----------------------------------------------------------------------

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *Api) handleWebsocketRequest() servemux.HandlerFunc[any] {
	return func(ctx *servemux.HttpCtx[any]) {
		// Check token
		apiKeys := ctx.GetHeader("Api-Key")
		if len(apiKeys) != 1 {
			ctx.Error(appError.NewErrUnauthorized("no api key present in API-Key header"))
			return
		}
		deviceIdentity, err := a.logic.GetDeviceIdentityFromToken(ctx.Context(), apiKeys[0])
		if err != nil {
			ctx.Error(appError.NewErrUnauthorized("token validation failed - %v", err))
			return
		}
		// Upgrade Websocke connection
		conn, err := upgrader.Upgrade(ctx.GetRawResponseWriter(), ctx.GetRawRequest(), nil)
		if err != nil {
			a.logger.Errorf("failed to upgrade websocket - %v", err)
			return
		}
		// Close connection
		defer func() {
			// Unregister connection
			a.logic.UnRegisterDevice(ctx.Context(), deviceIdentity.Id)
			// Send close message to client
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server closing connection")
			conn.WriteMessage(websocket.CloseMessage, closeMessage)
			// Close the connection
			conn.Close()
		}()
		// Register Device
		wsCtx, err := a.logic.RegisterDevice(ctx.Context(), deviceIdentity.Id, conn)
		if err != nil {
			a.logger.Errorf("failed to register device - %v", err)
			return
		}
		// Handle requests on Websocket interface
		<-wsCtx.Done()
	}
}

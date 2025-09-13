package ws

import (
	"context"
	"encoding/json"
	"sync"

	appError "github.com/fireops-software/fireops-core-edge-manager/error"
	"github.com/gorilla/websocket"
	"github.com/uoul/go-common/async"
)

type WsRequestClient struct {
	conn            *websocket.Conn
	ctx             context.Context
	stop            context.CancelFunc
	pendingRequests map[string]chan async.ActionResult[IResponse]
	requestCh       chan IRequest
	mux             *sync.RWMutex
}

func (w *WsRequestClient) Send(req IRequest) chan async.ActionResult[IResponse] {
	// Create Response Channel
	respCh := make(chan async.ActionResult[IResponse], 1)
	// Check if context is already cancelled
	select {
	case <-w.ctx.Done():
		respCh <- async.NewErrorActionResult[IResponse](
			appError.NewErrChannelClosed("client is shutting down"),
		)
		close(respCh)
		return respCh
	default:
	}
	w.mux.Lock()
	if _, exists := w.pendingRequests[req.GetMsgId()]; exists {
		w.mux.Unlock()
		respCh <- async.NewErrorActionResult[IResponse](
			appError.NewErrConflict("request with id %s already exists", req.GetMsgId()),
		)
		close(respCh)
		return respCh
	}
	w.pendingRequests[req.GetMsgId()] = respCh
	w.mux.Unlock()
	// Safe channel write with context check
	select {
	case w.requestCh <- req:
		// Success
	case <-w.ctx.Done():
		w.mux.Lock()
		delete(w.pendingRequests, req.GetMsgId())
		w.mux.Unlock()
		respCh <- async.NewErrorActionResult[IResponse](
			appError.NewErrChannelClosed("client is shutting down"),
		)
		close(respCh)
	default:
		// Channel full - could add timeout or make buffered
		w.mux.Lock()
		delete(w.pendingRequests, req.GetMsgId())
		w.mux.Unlock()
		respCh <- async.NewErrorActionResult[IResponse](
			appError.NewErrChannelClosed("request channel full"),
		)
		close(respCh)
	}
	// Return response channel
	return respCh
}

func (w *WsRequestClient) IsConnected() bool {
	select {
	case <-w.ctx.Done():
		return false
	default:
		return true
	}
}

func (w *WsRequestClient) readWs() {
	for {
		select {
		case <-w.ctx.Done():
			w.replyAll(async.NewErrorActionResult[IResponse](
				appError.NewErrChannelClosed("context canceled"),
			))
			return
		default:
			// Read Response
			resp := Response[json.RawMessage]{}
			if err := w.conn.ReadJSON(&resp); err != nil {
				w.stop()
				w.replyAll(async.NewErrorActionResult[IResponse](
					appError.NewErrNetwork("failed to read from websocket - %v", err),
				))
				return
			}
			// Find request for incomming data
			w.mux.Lock()
			respCh, exists := w.pendingRequests[resp.MsgId]
			if exists {
				respCh <- async.ActionResult[IResponse]{
					Result: &resp,
					Error:  nil,
				}
				close(respCh)
				delete(w.pendingRequests, resp.MsgId)
			}
			w.mux.Unlock()
		}
	}
}

func (w *WsRequestClient) Close() error {
	w.stop()
	return nil
}

func (w *WsRequestClient) Context() context.Context {
	return w.ctx
}

func (w *WsRequestClient) writeWs() {
	for {
		select {
		case req := <-w.requestCh:
			if err := w.conn.WriteJSON(req); err != nil {
				w.stop()
				w.replyAll(async.NewErrorActionResult[IResponse](
					appError.NewErrNetwork("failed to write to websocket - %v", err),
				))
				return
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *WsRequestClient) replyAll(resp async.ActionResult[IResponse]) {
	w.mux.Lock()
	defer w.mux.Unlock()
	for id, respCh := range w.pendingRequests {
		respCh <- resp
		close(respCh)
		delete(w.pendingRequests, id)
	}
}

func NewWsRequestClient(appCtx context.Context, conn *websocket.Conn) *WsRequestClient {
	ctx, cancel := context.WithCancel(appCtx)
	client := &WsRequestClient{
		conn:            conn,
		ctx:             ctx,
		stop:            cancel,
		pendingRequests: map[string]chan async.ActionResult[IResponse]{},
		requestCh:       make(chan IRequest, 50),
		mux:             &sync.RWMutex{},
	}
	go client.readWs()
	go client.writeWs()
	return client
}

package ws

import (
	"encoding/json"

	appError "github.com/fireops-software/fireops-core-edge-manager/error"
)

type Response[T any] struct {
	MsgId   string
	MsgType string
	Error   *string
	Body    T
}

// GetBody implements IResponse.
func (r *Response[T]) GetBody() json.RawMessage {
	data, _ := json.Marshal(r.Body)
	return data
}

// GetError implements IResponse.
func (r *Response[T]) GetError() error {
	if r.Error != nil {
		return appError.NewErrAgent("%s", *r.Error)
	}
	return nil
}

// GetMsgId implements IResponse.
func (r *Response[T]) GetMsgId() string {
	return r.MsgId
}

// GetMsgType implements IResponse.
func (r *Response[T]) GetMsgType() string {
	return r.MsgType
}

func NewResponse[T any](msgId string, msgType string, body T, err error) IResponse {
	var errStr *string
	if err != nil {
		str := err.Error()
		errStr = &str
	}
	return &Response[T]{
		MsgId:   msgId,
		MsgType: msgType,
		Body:    body,
		Error:   errStr,
	}
}

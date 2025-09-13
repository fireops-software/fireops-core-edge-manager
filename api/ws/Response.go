package ws

import "encoding/json"

type Response[T any] struct {
	MsgId   string
	MsgType string
	Error   error
	Body    T
}

// GetBody implements IResponse.
func (r *Response[T]) GetBody() json.RawMessage {
	data, _ := json.Marshal(r.Body)
	return data
}

// GetError implements IResponse.
func (r *Response[T]) GetError() error {
	return r.Error
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
	return &Response[T]{
		MsgId:   msgId,
		MsgType: msgType,
		Body:    body,
		Error:   err,
	}
}

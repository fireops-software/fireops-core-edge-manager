package ws

import "encoding/json"

type Request[T any] struct {
	MsgId   string
	MsgType string
	Body    T
}

// GetBody implements IRequest.
func (r *Request[T]) GetBody() json.RawMessage {
	data, _ := json.Marshal(r.Body)
	return data
}

// GetMsgId implements IRequest.
func (r *Request[T]) GetMsgId() string {
	return r.MsgId
}

// GetMsgType implements IRequest.
func (r *Request[T]) GetMsgType() string {
	return r.MsgType
}

func NewRequest[T any](msgId string, msgType string, body T) IRequest {
	return &Request[T]{
		MsgId:   msgId,
		MsgType: msgType,
		Body:    body,
	}
}

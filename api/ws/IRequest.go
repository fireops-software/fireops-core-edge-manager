package ws

import "encoding/json"

type IRequest interface {
	GetMsgId() string
	GetMsgType() string
	GetBody() json.RawMessage
}

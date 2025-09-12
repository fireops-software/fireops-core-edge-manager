package ws

type IRequest interface {
	GetMsgId() string
	GetMsgType() string
}

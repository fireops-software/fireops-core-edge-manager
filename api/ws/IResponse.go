package ws

type IResponse interface {
	IRequest
	GetError() error
}

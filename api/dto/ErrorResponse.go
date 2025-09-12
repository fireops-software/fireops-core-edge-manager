package dto

import "reflect"

type ErrorResponse struct {
	Type string
	Msg  string
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Type: reflect.TypeOf(err).Name(),
		Msg:  err.Error(),
	}
}

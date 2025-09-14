package api

import (
	"net/http"

	"github.com/fireops-software/fireops-core-edge-manager/api/dto"
	appError "github.com/fireops-software/fireops-core-edge-manager/error"
	"github.com/uoul/go-common/collections"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/servemux"
)

const (
	HEADER_API_KEY = "Api-Key"
)

func useErrorTranslation[T any](logger log.ILogger) servemux.HandlerFunc[T] {
	return func(ctx *servemux.HttpCtx[T]) {
		for _, err := range ctx.Errors() {
			logger.Error(err.Error())
			switch err.(type) {
			case appError.ErrNotFound:
				ctx.AbortWithResponse(http.StatusNotFound, dto.NewErrorResponse(err))
			case appError.ErrUnauthorized:
				ctx.AbortWithResponse(http.StatusUnauthorized, dto.NewErrorResponse(err))
			case appError.ErrConflict:
				ctx.AbortWithResponse(http.StatusConflict, dto.NewErrorResponse(err))
			case appError.ErrChannelClosed, appError.ErrNetwork:
				ctx.AbortWithResponse(http.StatusServiceUnavailable, dto.NewErrorResponse(err))
			case appError.ErrNotImplemented:
				ctx.AbortWithResponse(http.StatusNotImplemented, dto.NewErrorResponse(err))
			default:
				ctx.AbortWithResponse(http.StatusInternalServerError, dto.NewErrorResponse(err))
			}
		}
	}
}

func useApiKey[T any](validKeys []string) servemux.HandlerFunc[T] {
	return func(ctx *servemux.HttpCtx[T]) {
		keys := ctx.GetHeader(HEADER_API_KEY)
		// Check if API-Key header present
		if len(keys) <= 0 {
			ctx.AbortWithResponse(http.StatusUnauthorized, dto.NewErrorResponse(appError.NewErrUnauthorized("no API-Key header provided")))
			return
		}
		// Check if just one API-Key header
		if len(keys) != 1 {
			ctx.AbortWithResponse(http.StatusBadRequest, dto.NewErrorResponse(appError.ErrUnauthorized("only one API-Key header is allowed")))
			return
		}
		// Check if API-Key is valid
		if !collections.ContainsSlice(validKeys, func(k string) bool { return k == keys[0] }) {
			ctx.AbortWithResponse(http.StatusUnauthorized, dto.NewErrorResponse(appError.ErrUnauthorized("api key in API-Key header is not valid")))
			return
		}
	}
}

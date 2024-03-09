package middleware

import (
	"context"
	"net/http"

	"tomasweigenast.com/url-shortener/utils"
)

const (
	ContextRequestIdKey = "RequestId"
	HeaderNameRequestId = "X-Request-Id"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := utils.RequestID()
		w.Header().Set(HeaderNameRequestId, requestId)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextRequestIdKey, requestId)))
	})
}

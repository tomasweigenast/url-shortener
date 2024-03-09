package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Context().Value(ContextRequestIdKey).(string)
		log.Printf("Incoming request. Request Id [%s] Method [%s] From [%s]\n", requestId, r.Method, r.RemoteAddr)
		now := time.Now()
		next.ServeHTTP(w, r)
		diff := time.Now().Sub(now)
		log.Printf("Request finished. It took %d ms.\n", diff.Milliseconds())
	})
}

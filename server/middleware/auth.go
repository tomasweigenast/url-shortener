package middleware

import (
	"context"
	"net/http"
	"strings"

	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/server/response"
	"tomasweigenast.com/url-shortener/services"
)

const (
	HeaderNameAuthorization = "authorization"
	ContextUserKey          = "user"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get(HeaderNameAuthorization)
		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// invalid token, return earlier
			response.Failed(w, models.StringError{Reason: "unauthenticated"}, http.StatusUnauthorized)
			return
		}

		// validate token first
		token := parts[1]
		claims, err := services.AuthService().ValidateToken(r.Context(), token)
		if err != nil {
			response.Failed(w, models.StringError{Reason: "unauthenticated"}, http.StatusUnauthorized)
			return
		}

		// update context with new claims
		r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, models.User{
			Id:    claims.Sub,
			Name:  claims.Name,
			Email: claims.Email,
		}))

		next.ServeHTTP(w, r)
	})
}

func GetUid(req *http.Request) uint32 {
	return req.Context().Value(ContextUserKey).(models.User).Id
}

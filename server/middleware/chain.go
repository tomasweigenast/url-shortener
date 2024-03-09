package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain chains multiple middlewares with a fancy syntax. They are chained in the order they are declared
// handler is the actual HTTP handler and middlewares the list of middlewares to chain
func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.Handler {
	finalHandler := http.Handler(handler)
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}

	return finalHandler
}

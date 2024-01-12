package middleware

import (
	"WhisperWave-BackEnd/models"
	server "WhisperWave-BackEnd/server"
	"context"
	"net/http"
)

// Middleware function to inject the server object into the request context.
func WithServer(next http.Handler, server *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		srvCtx := models.ServerContext{Key: "server"}

		ctx = context.WithValue(ctx, srvCtx, server)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	server "WhisperWave-BackEnd/server"
	"WhisperWave-BackEnd/src/utils"
	"fmt"
	"net/http"
)

// middleware to check for authentication before accessing resources
func CheckAuthenticated(next http.Handler, server *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		if authToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing Authorization Header")
			return
		}

		// extract the actual token ID
		authToken = authToken[len("Bearer "):]

		err := utils.VerifyToken(authToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err)
			return
		}

		// if the user is authenticated, then server the requested resource
		next.ServeHTTP(w, r)
	})
}

func FormatValidator() {

}

func AddSecurityHeaders() {

}

package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Miskamyasa/utils/alerts"
	"github.com/Miskamyasa/utils/response"
)

// RecoveryMiddleware is a middleware that recovers from panics and sends an internal server error response
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Convert interface{} to an error
				var errMsg error
				if e, ok := err.(error); ok {
					errMsg = e
				} else {
					errMsg = fmt.Errorf("%v", err)
				}

				// Log the error and stack trace
				log.Printf("Recovered from panic: %v\nStack trace: %s", errMsg, debug.Stack())

				// Send alert and internal server error response
				alerts.Send("Panic recovery", errMsg)

				// Send internal server error response
				response.SendInternalServerError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks if the request has a valid auth token in the header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("auth-token")
		if token != os.Getenv("AUTH_TOKEN") {
			alerts.Send("Unauthorized request. Invalid auth token or token is nil", nil)
			response.SendInternalServerError(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

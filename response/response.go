package response

import (
	"encoding/json"
	"net/http"

	"github.com/Miskamyasa/utils/alerts"
)

type Response struct {
	message string
	payload interface{}
}

// NewResponse creates a new Response instance with type parameter T for payload
func NewResponse[T any](message string, payload T) Response {
	return Response{
		message: message,
		payload: payload,
	}
}

// SendJsonResponse sends a JSON response with the given payload and sets the Content-Type header to application/json.
func SendJsonResponse(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		alerts.Send("Error encoding the response", err)
		return
	}
}

// SendInternalServerError sends a 500 Internal Server Error response.
func SendInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("Internal Server Error"))
	if err != nil {
		alerts.Send("Error writing the response", err)
	}
}

// SendBadRequest sends a 400 Bad Request response with a custom message.
func SendBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Bad Request! " + msg))
	if err != nil {
		alerts.Send("Error writing the response", err)
	}
}

// HealthCheckHandler is a simple health check handler that responds with a 200 OK status.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		alerts.Send("Error writing the response", err)
		return
	}
}

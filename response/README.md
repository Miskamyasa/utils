# Package `response`

This package provides utility functions to handle HTTP responses in Go applications, including sending JSON data, handling internal server errors, bad requests, and health checks.

## Functions

### 1. SendJsonResponse
Sends a JSON response with the specified payload.

#### Parameters:
- `w http.ResponseWriter`: The response writer interface.
- `payload interface{}`: The payload to encode in the response body as JSON.

#### Example Usage:
```go
var data map[string]interface{} = map[string]interface{}{
    "key": "value",
}
SendJsonResponse(w, data)
```

### 2. SendInternalServerError
Sends a generic internal server error response (HTTP status code 500).

#### Parameters:
- `w http.ResponseWriter`: The response writer interface.

#### Example Usage:
```go
SendInternalServerError(w)
```

### 3. SendBadRequest
Sends a bad request response with an optional message (HTTP status code 400).

#### Parameters:
- `w http.ResponseWriter`: The response writer interface.
- `msg string`: Optional additional error message to include in the response body.

#### Example Usage:
```go
SendBadRequest(w, "Invalid input data")
```

### 4. HealthCheckHandler
Sends a health check response with an HTTP status code of 200 OK and a body containing "OK".

#### Parameters:
- `w http.ResponseWriter`: The response writer interface.
- `r *http.Request`: A pointer to the request information received from the client.

#### Example Usage:
```go
func someRouteHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/health" {
        HealthCheckHandler(w, r)
    }
}
```

## Notes:
- The package utilizes the `alert` function from the `github.com/Miskamyasa/utils/alerts` package to send alerts in case of errors.
- Ensure that you have the necessary imports and a running HTTP server to test these functions.

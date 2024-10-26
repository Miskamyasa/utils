# Package `middlewares`

This package provides various HTTP middleware functions to enhance the functionality and reliability of an HTTP server application.

## Functions

### 1. GenerateCacheKey
Generates a unique cache key based on the client's IP address and request path.

#### Parameters:
- `req *http.Request`: The incoming HTTP request.

#### Returns:
- A string representing the generated cache key.

#### Example Usage:
```go
cacheKey := GenerateCacheKey(req)
```

### 2. CacheMiddleware
A middleware function that checks if a response is cached and sends it from the cache instead of processing the request further if found in the cache.

#### Parameters:
- `next http.Handler`: The next handler in the chain to invoke after this middleware, which should handle the actual logic for generating a response.

#### Returns:
- An `http.Handler` that wraps around the original handler and adds caching functionality.

#### Example Usage:
```go
http.Handle("/", CacheMiddleware(http.HandlerFunc(someHandler)))
```

### 3. RecoveryMiddleware
A middleware function that recovers from panics, logs errors with stack traces, sends an alert if needed, and returns a generic internal server error response.

#### Parameters:
- `next http.Handler`: The next handler in the chain to invoke after this middleware.

#### Returns:
- An `http.Handler` that wraps around the original handler and adds recovery functionality for panics.

#### Example Usage:
```go
http.Handle("/", RecoveryMiddleware(http.HandlerFunc(someHandler)))
```

### 4. AuthMiddleware
A middleware function used for basic authentication, ensuring that incoming requests have a valid authorization token in their headers.

#### Parameters:
- `next http.Handler`: The next handler in the chain to invoke after this middleware.

#### Returns:
- An `http.Handler` that wraps around the original handler and adds basic authentication functionality.

#### Example Usage:
```go
http.Handle("/", AuthMiddleware(http.HandlerFunc(someHandler)))
```

## Notes:

- **CacheMiddleware**: This function uses a cache mechanism (imported from `github.com/Miskamyasa/utils/cache`) to check if the response for a given request is already cached. If found, it sends the cached content directly as a JSON response.
  
- **RecoveryMiddleware**: It catches any panics that occur in subsequent middleware or handlers and logs them along with stack traces. It also triggers an alert using `github.com/Miskamyasa/utils/alerts` package's `Send` function and returns an internal server error (HTTP status code 500) to the client.

- **AuthMiddleware**: This middleware checks for a specific authorization token in the request headers (`auth-token`). If the token does not match the configured `AUTH_TOKEN`, it logs an unauthorized access attempt and sends back an internal server error response.

## Requirements:

- Ensure that you have the necessary environment variables set up, such as the `AUTH_TOKEN` required by the `AuthMiddleware`.
  
- The middleware functions assume the existence of a caching system (`github.com/Miskamyasa/utils/cache`) and alerting utility (`github.com/Miskamyasa/utils/alerts`) packages.

---

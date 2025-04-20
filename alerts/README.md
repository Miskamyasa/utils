```markdown
# Package `alerts`

This package provides a simple and centralized logging and alerting mechanism for Go applications, built using the `zerolog` library. It aims to standardize logging output and provide functions for sending informational messages, errors, and fatal alerts that can terminate the application.

## Functions

### 1. `Send`
Logs a message with optional error details.

#### Parameters:
- `msg string`: The primary message to be logged.
- `err error`: An optional error object. If provided, the log will be recorded at the error level, otherwise at the info level.

#### Notes:
- This function creates a logger instance (or reuses an existing one) with the component field set to "alerts".
- If an error is provided, it logs the message and the error at the error level using `zerolog's Err` method.
- If no error is provided, it logs the message at the info level using `zerolog's Info` method.
- A `// TODO: send admin notification` comment indicates a placeholder for future functionality to send administrative alerts, which is not yet implemented.

#### Example Usage:
```go
package main

import (
	"fmt"
	"github.com/Miskamyasa/utils/alerts"
)

func main() {
	alerts.Send("Application started successfully", nil)

	err := fmt.Errorf("example error occurred")
	alerts.Send("Something went wrong", err)
}
```

### 2. `Fatal`
Logs a message with optional error details and then terminates the application.

#### Parameters:
- `msg string`: The primary message to be logged.
- `err error`: An optional error object. If provided, the log will be recorded with error details.

#### Notes:
- This function first calls the `Send` function to log the provided message and error (if any).
- After logging, it terminates the application immediately by calling `os.Exit(1)`. This is typically used for unrecoverable errors that prevent the application from continuing to run safely.

#### Example Usage:
```go
package main

import (
	"fmt"
	"github.com/Miskamyasa/utils/alerts"
)

func main() {
	configErr := fmt.Errorf("failed to load configuration")
	if configErr != nil {
		alerts.Fatal("Application cannot start due to configuration error", configErr)
		// The application will exit here after logging the fatal error.
	}
	fmt.Println("This line will not be reached if Fatal is called.")
}
```

### 3. `CreateLogger`
Creates and returns a `zerolog.Logger` instance.

#### Returns:
- `zerolog.Logger`: A configured `zerolog.Logger` instance.

#### Notes:
- This function uses `sync.Once` to ensure that the logger is initialized only once throughout the application's lifecycle. Subsequent calls will return the same logger instance.
- The logger is configured to write to `os.Stdout` using `zerolog.ConsoleWriter` with `time.RFC3339` time format for human-readable output in development.
- It includes the following fields in every log message:
    - `timestamp`: The current timestamp in RFC3339 format.
    - `service`: The service name, obtained from the `SERVICE_NAME` environment variable.
    - `version`: The service version, obtained from the `SERVICE_VERSION` environment variable.
    - `env`: The environment name, obtained from the `ENV` environment variable.
- **Environment Variables**: This function relies on the following environment variables for context in logs:
    - `SERVICE_NAME`:  The name of the service. This should be set to identify the source of the logs.
    - `SERVICE_VERSION`: The version of the service. Useful for tracking logs across different versions.
    - `ENV`: The environment the service is running in (e.g., "development", "production", "staging").

#### Example Usage:
```go
package main

import (
	"github.com/Miskamyasa/utils/alerts"
	"os"
)

func main() {
	os.Setenv("SERVICE_NAME", "my-app")
	os.Setenv("SERVICE_VERSION", "1.0.0")
	os.Setenv("ENV", "development")

	logger := alerts.CreateLogger()
	logger.Info().Msg("Application logger initialized")

	// You can use the logger directly for more customized logging:
	logger.Error().Str("component", "main").Msg("Example error log using direct logger")
}
```

## Notes:
- **Centralized Logging**: This package provides a central point for logging across your application, ensuring consistent formatting and contextual information in your logs.
- **Error Handling**: Use `alerts.Send` for general logging and error reporting. Use `alerts.Fatal` for critical errors that require application termination.
- **Environment Configuration**: Make sure to set the `SERVICE_NAME`, `SERVICE_VERSION`, and `ENV` environment variables to provide context to your logs.
- **Dependencies**: This package depends on `github.com/rs/zerolog`. Ensure this dependency is included in your `go.mod` file.
```
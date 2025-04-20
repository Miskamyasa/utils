```markdown
# Package `cache`

This package provides a caching mechanism using Redis and an optional in-memory TinyLFU cache for Go applications. It simplifies common caching operations such as initialization, key generation, setting, and retrieving data from the cache.

## Functions

### 1. `InitCache`
Initializes the Redis client and sets up the cache instance with an optional TinyLFU local cache.

#### Returns:
- `*redis.Client`: A pointer to the initialized Redis client.

#### Configuration:
This function relies on environment variables for configuration:
- `REDIS_URL`:  The connection string for the Redis server (e.g., `redis://localhost:6379`). This is **required**. If not set, the application will terminate with a fatal error.
- `LFU_SIZE`: (Optional) The size of the TinyLFU local cache. If not set or if parsing fails, it defaults to `1000`.

#### Notes:
- This function establishes a connection to Redis and initializes a global cache instance used by other functions in the package.
- It uses `alerts.Fatal` from `github.com/Miskamyasa/utils/alerts` to handle Redis connection errors, causing the application to exit if the connection fails.
- The TinyLFU cache acts as a first-level cache to improve performance by storing frequently accessed items in memory.

#### Example Usage:
```go
package main

import (
	"fmt"
	"github.com/Miskamyasa/utils/cache"
)

func main() {
	redisClient := cache.InitCache()
	fmt.Println("Redis client initialized:", redisClient != nil)
	// You can now use other cache functions in your application.
}
```

### 2. `GenerateCacheKey`
Generates a cache key string based on the incoming HTTP request.

#### Parameters:
- `req *http.Request`: The HTTP request object.

#### Returns:
- `string`: A generated cache key in the format `"cache:<IP address>:<URL path>"`.

#### Notes:
- This function uses the remote IP address and URL path from the HTTP request to create a unique cache key.
- It's useful for caching responses based on the request's origin and the requested resource.

#### Example Usage:
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/Miskamyasa/utils/cache"
)

func main() {
	req, _ := http.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "192.168.1.100:12345" // Example Remote Address
	key := cache.GenerateCacheKey(req)
	fmt.Println("Generated Cache Key:", key) // Output: Generated Cache Key: cache:192.168.1.100:12345:/api/data
}
```

### 3. `CreateDuration`
Creates a `time.Duration` from an integer representing seconds.

#### Parameters:
- `seconds int`: The number of seconds for the duration.

#### Returns:
- `time.Duration`: A `time.Duration` object representing the specified number of seconds.

#### Example Usage:
```go
package main

import (
	"fmt"
	"time"
	"github.com/Miskamyasa/utils/cache"
)

func main() {
	duration := cache.CreateDuration(60) // 60 seconds
	fmt.Println("Duration:", duration)     // Output: Duration: 1m0s
	fmt.Printf("Duration type: %T\n", duration) // Output: Duration type: time.Duration
}
```

### 4. `GetCache[T any]`
Retrieves data from the cache associated with the given key and unmarshals it into the provided payload.

#### Type Parameter:
- `T any`:  The type of the payload to be retrieved from the cache. This function is generic and can work with any data type.

#### Parameters:
- `key string`: The cache key to retrieve the data from.
- `payload *T`: A pointer to a variable of type `T` where the retrieved and unmarshaled data will be stored.

#### Returns:
- `error`: Returns an error if:
    - There is an issue retrieving data from the cache.
    - There is an issue unmarshaling the retrieved JSON data into the payload.
    - Returns `nil` if the data is successfully retrieved and unmarshaled, or if the environment is set to "development" (in which case, the cache is bypassed and no error is returned).

#### Notes:
- In "development" environment (determined by the `ENV` environment variable), this function bypasses the cache and returns `nil` immediately. This is useful for local development where caching might not be desired.
- Data is stored in the cache as JSON bytes. This function automatically unmarshals the JSON data into the provided payload.
- If the key is not found in the cache, `instance.Get` will return an error (likely `redis.Nil`), which will be propagated by this function.

#### Example Usage:
```go
package main

import (
	"fmt"
	"github.com/Miskamyasa/utils/cache"
	"os"
)

type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	os.Setenv("REDIS_URL", "redis://localhost:6379") // Set Redis URL for example
	redisClient := cache.InitCache()
	defer redisClient.Close()

	key := "user:123"
	var userData UserData

	err := cache.GetCache(key, &userData)
	if err != nil {
		fmt.Println("Error getting cache:", err)
		// Handle cache miss or other errors, e.g., fetch data from source and SetCache
	} else {
		fmt.Println("Retrieved from cache:", userData)
	}
}
```

### 5. `SetCache[T any]`
Stores data in the cache with the given key and Time-To-Live (TTL).

#### Type Parameter:
- `T any`: The type of the payload to be stored in the cache. This function is generic and can work with any data type.

#### Parameters:
- `key string`: The cache key to store the data under.
- `payload T`: The data payload to be stored in the cache. This will be marshaled to JSON before being stored.
- `TTL time.Duration`: The duration for which the data should be kept in the cache before expiring.

#### Returns:
- `error`: Returns an error if:
    - There is an issue marshaling the payload into JSON.
    - There is an issue setting the data in the cache.
    - Returns `nil` if the data is successfully set in the cache, or if the environment is set to "development" (in which case, the cache is bypassed and no error is returned).

#### Notes:
- In "development" environment (determined by the `ENV` environment variable), this function bypasses the cache and returns `nil` immediately.
- Data is stored in the cache as JSON bytes after marshaling the provided payload.

#### Example Usage:
```go
package main

import (
	"fmt"
	"time"
	"os"
	"github.com/Miskamyasa/utils/cache"
)

type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	os.Setenv("REDIS_URL", "redis://localhost:6379") // Set Redis URL for example
	redisClient := cache.InitCache()
	defer redisClient.Close()

	key := "user:123"
	userData := UserData{ID: 123, Name: "John Doe"}
	ttl := cache.CreateDuration(300) // 5 minutes

	err := cache.SetCache(key, userData, ttl)
	if err != nil {
		fmt.Println("Error setting cache:", err)
	} else {
		fmt.Println("Data set in cache successfully")
	}
}
```

## Notes:
- **Environment Variable `ENV`**:  If the environment variable `ENV` is set to `"development"`, both `GetCache` and `SetCache` functions will effectively bypass the cache operations. This is intended for development environments where actual caching might not be desired.
- **Error Handling**:  Functions in this package return errors to allow for proper error handling in the calling code. Always check the returned error and handle cache misses or failures appropriately in your application logic.
- **JSON Serialization**: Data is serialized to JSON before being stored in the cache and deserialized when retrieved. Ensure that the data structures you are caching are JSON-serializable.
- **Dependencies**: This package depends on `github.com/go-redis/cache/v9`, `github.com/redis/go-redis/v9`, and `github.com/Miskamyasa/utils/alerts`. Make sure to include these dependencies in your `go.mod` file.
```
package cache

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/Miskamyasa/utils/alerts"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var instance *cache.Cache

var cacheCtx = context.Background()

func InitCache() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	_, err := client.Ping(cacheCtx).Result()
	if err != nil {
		alerts.Fatal("Failed to connect to Redis", err)
	}

	LFUSize, err := strconv.Atoi(os.Getenv("LFU_SIZE"))
	if err != nil {
		LFUSize = 1000
	}

	instance = cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(LFUSize, time.Second),
	})

	return client
}

// CreateDuration creates a time.Duration from seconds.
func CreateDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// GetCache retrieves a value from the cache using the provided key and unmarshals it into the provided payload.
func GetCache[T any](key string, payload *T) error {
	if os.Getenv("ENV") == "development" {
		return nil
	}

	var jsonPayload []byte
	err := instance.Get(cacheCtx, key, &jsonPayload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonPayload, payload)
	if err != nil {
		return err
	}

	return nil
}

// SetCache stores a value in the cache with the provided key and TTL (time-to-live).
func SetCache[T any](key string, payload T, TTL time.Duration) error {
	if os.Getenv("ENV") == "development" {
		return nil
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = instance.Set(&cache.Item{
		Ctx:   cacheCtx,
		Key:   key,
		Value: jsonPayload, // Store JSON bytes
		TTL:   TTL,
	})
	if err != nil {
		return err
	}

	return nil
}

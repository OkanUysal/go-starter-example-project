package config

import (
	"os"
	"strconv"
	"time"

	"github.com/OkanUysal/go-cache"
	"github.com/OkanUysal/go-logger"
)

var Cache *cache.Cache

// InitCache initializes the cache based on environment configuration
func InitCache() error {
	cacheType := GetEnv("CACHE_TYPE", "memory")
	ttlSeconds := GetEnv("CACHE_TTL", "300")

	Logger.Info("Cache configuration", logger.String("type", cacheType), logger.String("ttl", ttlSeconds))

	ttl, err := strconv.Atoi(ttlSeconds)
	if err != nil || ttl <= 0 {
		Logger.Warn("Invalid TTL value, using default", logger.String("value", ttlSeconds), logger.Int("default", 300))
		ttl = 300 // default 5 minutes
	}

	config := &cache.Config{
		Backend:         cache.BackendMemory,
		DefaultTTL:      time.Duration(ttl) * time.Second,
		CleanupInterval: 10 * time.Minute, // Clean up expired entries every 10 minutes
	}

	// Set Redis backend if configured
	if cacheType == "redis" {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			Logger.Warn("REDIS_URL not set, falling back to memory cache")
			config.Backend = cache.BackendMemory
		} else {
			config.Backend = cache.BackendRedis
			config.RedisURL = redisURL
		}
	}

	Cache, err = cache.New(config)
	if err != nil {
		return err
	}

	Logger.Info("Cache initialized successfully",
		logger.String("backend", string(config.Backend)),
		logger.Int("ttl_seconds", ttl))

	return nil
}

// GetCache returns the cache instance
func GetCache() *cache.Cache {
	return Cache
}

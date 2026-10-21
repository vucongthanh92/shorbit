package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, RedisAddr, RedisPassword string
	RedisDB                                         int
	CacheTTL                                        time.Duration
}

func Load() Config {
	ttl, err := time.ParseDuration(value("CACHE_TTL", "5m"))
	if err != nil {
		ttl = 5 * time.Minute
	}
	redisDB, _ := strconv.Atoi(value("REDIS_DB", "0"))
	return Config{
		HTTPAddr:    value("HTTP_ADDR", ":8080"),
		DatabaseURL: value("DATABASE_URL", "postgres://shortbit:shortbit@localhost:5432/shortbit?sslmode=disable"),
		RedisAddr:   value("REDIS_ADDR", "localhost:6379"), RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB: redisDB, CacheTTL: ttl,
	}
}
func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

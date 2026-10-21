package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(addr, password string, db int, ttl time.Duration) *Redis {
	return &Redis{client: redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db}), ttl: ttl}
}
func (r *Redis) GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	raw, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(raw, dst)
}
func (r *Redis) SetJSON(ctx context.Context, key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, raw, r.ttl).Err()
}
func (r *Redis) Delete(ctx context.Context, key string) error { return r.client.Del(ctx, key).Err() }
func (r *Redis) Close() error                                 { return r.client.Close() }

package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// NewRedis creates a new Redis client connection
func NewRedis(addr, password string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Verify connection with ping
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

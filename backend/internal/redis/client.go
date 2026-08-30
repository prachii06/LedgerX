package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/prachii06/LedgerX/internal/config"
)

type Client struct {
	*redis.Client
}

func NewClient(cfg *config.Config) *Client {
	addr := fmt.Sprintf(
		"%s:%s",
		cfg.RedisHost,
		cfg.RedisPort,
	)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Client{
		Client: client,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
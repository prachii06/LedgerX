package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/prachii06/LedgerX/internal/models"
)

type ReconciliationCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewReconciliationCache(
	client *redis.Client,
	ttl time.Duration,
) *ReconciliationCache {
	return &ReconciliationCache{
		client: client,
		ttl:    ttl,
	}
}


//to check whether reconciliation result exists in Redis
func (c *ReconciliationCache) Get(
	ctx context.Context,
	transactionID string,
) (*models.ReconciliationResult, error) {

	key := fmt.Sprintf(
		"reconciliation:%s",
		transactionID,
	)

	data, err := c.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var result models.ReconciliationResult

	if err := json.Unmarshal(
		[]byte(data),
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}


//to store a new reconciliation result
func (c *ReconciliationCache) Set(
	ctx context.Context,
	result *models.ReconciliationResult,
) error {

	key := fmt.Sprintf(
		"reconciliation:%s",
		result.TransactionID,
	)

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		key,
		data,
		c.ttl,
	).Err()
}


//to remove stale cached result
func (c *ReconciliationCache) Delete(
	ctx context.Context,
	transactionID string,
) error {

	key := fmt.Sprintf(
		"reconciliation:%s",
		transactionID,
	)

	return c.client.Del(ctx, key).Err()
}
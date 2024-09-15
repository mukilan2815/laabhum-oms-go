package cache

import (
	"golang.org/x/net/context"
)

// Removed duplicate NewRedisCache, Get, and Set methods

func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

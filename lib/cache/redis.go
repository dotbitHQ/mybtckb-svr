package cache

import (
	"context"
	"github.com/go-redis/redis"
)

type RedisCache struct {
	Ctx context.Context
	Red *redis.Client
}

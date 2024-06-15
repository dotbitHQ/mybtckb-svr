package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type RedisCache struct {
	Ctx context.Context
	Red *redis.Client
}

func BtcAddrKey(txHash string, index uint32) string {
	return fmt.Sprintf("btc:addr:%s-%d", txHash, index)
}
func (r *RedisCache) GetBtcAddrCache(txHash string, index uint32) (string, error) {
	key := strings.ToLower(BtcAddrKey(txHash, index))
	fmt.Println("redis get key: ", key)
	if value, err := r.Red.Get(key).Result(); err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("error querying key %s err: %s", key, err)
	} else {
		return value, nil
	}

}
func (r *RedisCache) SetBtcAddrCache(txHash, addr string, index uint32) error {
	if r.Red == nil {
		return fmt.Errorf("redis is nil")
	}
	key := BtcAddrKey(txHash, index)
	return r.Red.Set(key, addr, time.Hour*24).Err()
}

package cache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	m map[string]*cacheEntry
	sync.RWMutex
}

type cacheEntry struct {
	key       string
	value     interface{}
	timestamp int64
}

func NewCache(ctx context.Context, wg *sync.WaitGroup) *Cache {
	c := &Cache{
		m: make(map[string]*cacheEntry),
	}
	ticker := time.NewTicker(time.Minute)
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				for k, v := range c.m {
					if time.Now().Unix() >= v.timestamp {
						c.Delete(k)
					}
				}
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()
	return c
}

func (c *Cache) Set(key string, value interface{}, expires ...time.Duration) {
	c.Lock()
	defer c.Unlock()

	var expire int64
	if len(expires) > 0 {
		expire = time.Now().Unix() + int64(expires[0].Seconds())
	}

	c.m[key] = &cacheEntry{
		key:       key,
		value:     value,
		timestamp: expire,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	val, ok := c.m[key]
	if !ok {
		c.RUnlock()
		return nil, false
	}
	if time.Now().Unix() >= val.timestamp {
		c.RUnlock()
		c.Lock()
		delete(c.m, key)
		c.Unlock()
		return nil, false
	}
	c.RUnlock()
	return val.value, true
}

func (c *Cache) Delete(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.m[key]
	if !ok {
		return nil, false
	}
	delete(c.m, key)

	if time.Now().Unix() >= val.timestamp {
		return false, false
	}
	return val.value, true
}

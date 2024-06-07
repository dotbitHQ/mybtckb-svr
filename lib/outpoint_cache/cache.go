package outpoint_cache

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Cache struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	rw          sync.RWMutex
	mapOutPoint map[string]int64
}

func NewCache(ctx context.Context, wg *sync.WaitGroup) *Cache {
	return &Cache{
		ctx:         ctx,
		wg:          wg,
		rw:          sync.RWMutex{},
		mapOutPoint: make(map[string]int64),
	}
}

func (d *Cache) AddOutPoint(outPoint []string) {
	if len(outPoint) == 0 {
		return
	}
	d.rw.Lock()
	defer d.rw.Unlock()
	for _, v := range outPoint {
		d.mapOutPoint[v] = time.Now().Unix()
	}
}

func (d *Cache) clearExpiredOutPoint(t time.Duration) {
	d.rw.Lock()
	defer d.rw.Unlock()
	timestamp := time.Now().Add(-t).Unix()
	before := len(d.mapOutPoint)
	for k, v := range d.mapOutPoint {
		if v < timestamp {
			delete(d.mapOutPoint, k)
		}
	}
	log.WithFields(log.Fields{
		"before": before,
		"after":  len(d.mapOutPoint),
	}).Info("clearExpiredOutPoint")
}

func (d *Cache) ExistOutPoint(outPoint string) bool {
	d.rw.RLock()
	defer d.rw.RUnlock()
	if _, ok := d.mapOutPoint[outPoint]; ok {
		return true
	}
	return false
}

func (d *Cache) RunClearExpiredOutPoint(t time.Duration) {
	ticker := time.NewTicker(t)
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.clearExpiredOutPoint(t)
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

package server

import (
	"sync"
	"time"
)

type Cache struct {
	data sync.Map
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.data.Store(key, value)

	go func() {
		<-time.After(ttl)
		c.data.Delete(key)
	}()
}

func (c *Cache) Get(key string) (any, bool) {
	return c.data.Load(key)
}

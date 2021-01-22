package cache

import (
	"errors"
	"sync"
)

type MemoryCache struct {
	cache map[interface{}]interface{}
	sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[interface{}]interface{}),
	}
}

func (c *MemoryCache) Set(key, value interface{}) error {
	c.Lock()
	c.cache[key] = value
	c.Unlock()

	return nil
}

func (c *MemoryCache) Get(key interface{}) (interface{}, error) {
	c.RLock()
	value, ex := c.cache[key]
	c.RUnlock()

	if !ex {
		return nil, errors.New("not found")
	}

	return value, nil
}

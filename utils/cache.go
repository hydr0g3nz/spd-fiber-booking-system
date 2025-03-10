package utils

import (
	"sync"
)

// Cache defines the interface for cache operations
type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	GetAll() map[string]interface{}
}

// InMemoryCache implements Cache interface with in-memory storage
type InMemoryCache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.data[key]
	return val, exists
}

func (c *InMemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *InMemoryCache) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]interface{}),
	}
}

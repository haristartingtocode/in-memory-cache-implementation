package cache_generics

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem[T any] struct {
	Value     T
	ExpiresAt time.Time
}

type Cache[K comparable, T any] struct {
	Mu   sync.RWMutex
	Data map[K]CacheItem[T]
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	newCache := &Cache[K, T]{
		Data: make(map[K]CacheItem[T]),
	}

	go newCache.CleanExpiredCacheItems()

	return newCache
}

func (c *Cache[K, T]) Set(key K, value T, timeToExpire time.Duration) {
	cacheExpireTime := time.Now().Add(timeToExpire)
	cacheItem := CacheItem[T]{Value: value, ExpiresAt: cacheExpireTime}
	c.Mu.Lock()
	c.Data[key] = cacheItem
	c.Mu.Unlock()
}

func (c *Cache[K, T]) Get(key K) (value T, ok bool) {
	c.Mu.RLock()
	CacheItem, ok := c.Data[key]
	c.Mu.RUnlock()
	if !ok {
		var zero T
		return zero, false
	}
	return CacheItem.Value, ok
}

func (c *Cache[K, T]) Count() int {
	c.Mu.RLock()
	dataLen := len(c.Data)
	c.Mu.RUnlock()
	return dataLen
}

func (c *Cache[K, T]) PrintStruct() {
	c.Mu.RLock()
	fmt.Printf("map is %v\n", c.Data)
	c.Mu.RUnlock()
}

func (c *Cache[K, T]) CleanExpiredCacheItems() {
	for {
		fmt.Println("Waking up and cleaning the expired caches")
		currentTime := time.Now()
		c.Mu.Lock()
		for key, value := range c.Data {
			if currentTime.After(value.ExpiresAt) {
				delete(c.Data, key)
			}
		}
		c.Mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func (c *Cache[K, T]) Delete(key K) {
	c.Mu.Lock()
	delete(c.Data, key)
	c.Mu.Unlock()
}

func (c *Cache[K, T]) Clear() {
	c.Mu.Lock()
	for key := range c.Data {
		delete(c.Data, key)
	}
	c.Mu.Unlock()
}

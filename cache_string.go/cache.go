package cache_string

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Value     string
	ExpiresAt time.Time
}

type Cache struct {
	Mu   sync.RWMutex
	Data map[string]*CacheItem
}

func NewCache() *Cache {
	newCache := &Cache{
		Data: make(map[string]*CacheItem),
	}

	go newCache.CleanExpiredCacheItems()

	return newCache
}

func (c *Cache) Set(key, value string, timeToExpire time.Duration) {
	cacheExpireTime := time.Now().Add(timeToExpire)
	cacheItem := &CacheItem{Value: value, ExpiresAt: cacheExpireTime}
	c.Mu.Lock()
	c.Data[key] = cacheItem
	c.Mu.Unlock()
}

func (c *Cache) Get(key string) (value string, ok bool) {
	c.Mu.RLock()
	CacheItem, ok := c.Data[key]
	c.Mu.RUnlock()
	if !ok {
		return "", false
	}
	return CacheItem.Value, ok
}

func (c *Cache) Count() int {
	c.Mu.RLock()
	dataLen := len(c.Data)
	c.Mu.RUnlock()
	return dataLen
}

func (c *Cache) PrintStruct() {
	c.Mu.RLock()
	fmt.Printf("map is %v\n", c.Data)
	c.Mu.RUnlock()
}

func (c *Cache) CleanExpiredCacheItems() {
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

func (c *Cache) Delete(key string) {
	c.Mu.Lock()
	delete(c.Data, key)
	c.Mu.Unlock()
}

func (c *Cache) Clear() {
	c.Mu.Lock()
	for key := range c.Data {
		delete(c.Data, key)
	}
	c.Mu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(5)

	cache := NewCache()

	for range 5 {
		go func() {
			fmt.Println("Value to be set")
			cache.Set("a", "b", 3*time.Second)
			fmt.Println("Count", cache.Count())
			cache.PrintStruct()
			wg.Done()
		}()
	}

	// for range 5 {
	// 	go func() {
	// 		fmt.Println("Going to sleep")
	// 		time.Sleep(10 * time.Second)
	// 		fmt.Println("Waking up")

	// 		value, ok := cache.Get("a")
	// 		fmt.Println("ok", ok)
	// 		if ok {
	// 			fmt.Println("a", value)
	// 		} else {
	// 			fmt.Println("not present")
	// 		}
	// 		wg.Done()
	// 	}()
	// }

	cache.Clear()

	value, ok := cache.Get("a")
	if ok {
		fmt.Println("a", value)
	} else {
		fmt.Println("not present")
	}

	wg.Wait()

}

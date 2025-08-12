package main

import (
	"fmt"
	"sync"
	"time"

	"in-memory-cache-implementation/cache_generics"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)

	cache := cache_generics.NewCache[string, string]()

	for range 5 {
		go func() {
			fmt.Println("Value to be set")
			cache.Set("a", "b", 3*time.Second)
			fmt.Println("Count", cache.Count())
			cache.PrintStruct()
			wg.Done()
		}()
	}

	for range 5 {
		go func() {
			fmt.Println("Going to sleep")
			time.Sleep(10 * time.Second)
			fmt.Println("Waking up")

			value, ok := cache.Get("a")
			fmt.Println("ok", ok)
			if ok {
				fmt.Println("a", value)
			} else {
				fmt.Println("not present")
			}
			wg.Done()
		}()
	}

	cache.Clear()

	value, ok := cache.Get("a")
	if ok {
		fmt.Println("a", value)
	} else {
		fmt.Println("not present")
	}

	wg.Wait()

}

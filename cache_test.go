package main

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()

	if cache == nil {
		t.Fatal("NewCache() returned nil")
	}

	if cache.Data == nil {
		t.Error("Cache.Data map is nil")
	}

	if len(cache.Data) != 0 {
		t.Errorf("Expected empty cache, got %d items", len(cache.Data))
	}
}

func TestSet(t *testing.T) {
	cache := NewCache()
	cache.Set("a", "b", 5)

	cacheItem, ok := cache.Data["a"]

	if !ok {
		t.Fatal("Error in Set")
	}

	if cacheItem.Value != "b" {
		t.Fatal("Error in Set")
	}
}

func TestGet(t *testing.T) {
	cache := NewCache()

	cache.Set("a", "b", 5)

	value, ok := cache.Get("a")

	if !ok {
		t.Fatal("Error in Set")
	}

	if value != "b" {
		t.Fatal("Error in Set")
	}
}

func TestDelete(t *testing.T) {
	cache := NewCache()

	cache.Set("a", "b", 5)

	cache.Delete("a")

	_, ok := cache.Get("a")

	if ok {
		t.Fatal("Error in Set")
	}
}

func TestClear(t *testing.T) {
	cache := NewCache()

	cache.Set("a", "b", 5)

	cache.Clear()

	_, ok := cache.Get("a")

	if ok {
		t.Fatal("Error in Set")
	}
}

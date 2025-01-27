package cache

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHttpCache_BasicOperations(t *testing.T) {
	cache := NewHttpCache()

	// Test Add and Get
	url := "http://example.com"
	data := []byte("test data")
	cache.Add(url, data)

	value, exists := cache.Get(url)
	if !exists {
		t.Error("Expected cache entry to exist")
	}
	if string(value) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(value))
	}

	// Test Delete
	cache.Delete(url)
	_, exists = cache.Get(url)
	if exists {
		t.Error("Expected cache entry to be deleted")
	}
}

func TestHttpCache_ConcurrentAccess(t *testing.T) {
	cache := NewHttpCache()
	var wg sync.WaitGroup
	iterations := 100

	// Concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := "http://example.com"
			data := []byte("test data")
			cache.Add(url, data)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := "http://example.com"
			cache.Get(url)
		}(i)
	}

	wg.Wait()
}

func TestCachedHttpClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := NewCachedHttpClient()

	// First request should hit the server
	data1, err := client.GetWithCache(server.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if string(data1) != "test response" {
		t.Errorf("Expected 'test response', got '%s'", string(data1))
	}

	// Second request should hit the cache
	data2, err := client.GetWithCache(server.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if string(data2) != "test response" {
		t.Errorf("Expected 'test response', got '%s'", string(data2))
	}

	// Verify both responses are identical (from cache)
	if string(data1) != string(data2) {
		t.Error("Expected cached response to be identical to original response")
	}
}

func TestHttpCache_Expiration(t *testing.T) {
	// Create a manual ticker for controlled testing
	reapC := make(chan time.Time)
	done := make(chan struct{}) // channel for synchronization
	cache := &HttpCache{
		cache: make(map[string]CacheEntry),
		reapC: reapC,
	}

	// Add test entries with different creation times
	cache.mu.Lock()
	cache.cache["fresh"] = CacheEntry{
		Value:     []byte("fresh data"),
		CreatedAt: time.Now(),
	}
	cache.cache["stale"] = CacheEntry{
		Value:     []byte("stale data"),
		CreatedAt: time.Now().Add(-2 * time.Hour), // 2 hours old, should be reaped
	}
	cache.mu.Unlock()

	// Start the reap loop with synchronization
	go func() {
		cache.reapLoop()
		close(done)
	}()

	// Trigger the reaping
	reapC <- time.Now()
	close(reapC) // Close reapC to exit reapLoop
	<-done       // Wait for reapLoop to finish

	// Verify that only the fresh entry remains
	cache.mu.Lock()
	if len(cache.cache) != 1 {
		t.Errorf("Expected 1 entry in cache, got %d", len(cache.cache))
	}
	if _, exists := cache.cache["stale"]; exists {
		t.Error("Stale entry should have been reaped")
	}
	if _, exists := cache.cache["fresh"]; !exists {
		t.Error("Fresh entry should still exist")
	}
	cache.mu.Unlock()
}

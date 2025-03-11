package utils_test

import (
	"sync"
	"testing"

	"github.com/hydr0g3nz/spd-fiber-booking-system/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryCache(t *testing.T) {
	// Act
	cache := utils.NewInMemoryCache()

	// Assert
	assert.NotNil(t, cache, "Expected cache instance to be non-nil")
}

func TestSetAndGet(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	key := "test-key"
	value := "test-value"

	// Act
	cache.Set(key, value)
	retrievedValue, exists := cache.Get(key)

	// Assert
	assert.True(t, exists, "Expected key to exist in cache")
	assert.Equal(t, value, retrievedValue, "Retrieved value doesn't match what was stored")
}

func TestGetNonExistentKey(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	key := "non-existent-key"

	// Act
	value, exists := cache.Get(key)

	// Assert
	assert.False(t, exists, "Expected key to not exist in cache")
	assert.Nil(t, value, "Expected nil value for non-existent key")
}

func TestDelete(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	key := "test-key"
	value := "test-value"
	cache.Set(key, value)

	// Verify key exists before deletion
	_, exists := cache.Get(key)
	assert.True(t, exists, "Expected key to exist before deletion")

	// Act
	cache.Delete(key)

	// Assert
	retrievedValue, exists := cache.Get(key)
	assert.False(t, exists, "Expected key to not exist after deletion")
	assert.Nil(t, retrievedValue, "Expected nil value after deletion")
}

func TestGetAll(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Act
	allData := cache.GetAll()

	// Assert
	assert.Equal(t, 3, len(allData), "Expected 3 items in cache")
	assert.Equal(t, "value1", allData["key1"])
	assert.Equal(t, "value2", allData["key2"])
	assert.Equal(t, "value3", allData["key3"])
}

func TestConcurrentAccess(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	const goroutines = 10
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // For both readers and writers

	// Act - Test concurrent writes
	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "concurrent-key-" + string(rune('A'+index))
				cache.Set(key, j)
			}
		}(i)
	}

	// Act - Test concurrent reads
	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "concurrent-key-" + string(rune('A'+index))
				cache.Get(key)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// If we got here without deadlocks or race conditions, the test passes
	// We can also verify that we have the expected number of keys
	allData := cache.GetAll()
	assert.Equal(t, goroutines, len(allData), "Expected number of keys doesn't match the number of writer goroutines")
}

func TestSetOverwritesExistingValue(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	key := "test-key"
	initialValue := "initial-value"
	updatedValue := "updated-value"

	// Act
	cache.Set(key, initialValue)
	cache.Set(key, updatedValue) // Overwrite
	retrievedValue, exists := cache.Get(key)

	// Assert
	assert.True(t, exists, "Expected key to exist in cache")
	assert.Equal(t, updatedValue, retrievedValue, "Retrieved value should be the updated value")
}

func TestCachePersistsComplexTypes(t *testing.T) {
	// Arrange
	cache := utils.NewInMemoryCache()
	key := "object-key"

	// Complex struct to store
	type TestStruct struct {
		ID   int
		Name string
		Tags []string
	}

	value := TestStruct{
		ID:   123,
		Name: "Test Object",
		Tags: []string{"tag1", "tag2", "tag3"},
	}

	// Act
	cache.Set(key, value)
	retrievedValue, exists := cache.Get(key)

	// Assert
	assert.True(t, exists, "Expected key to exist in cache")

	// Type assertion to convert interface{} back to TestStruct
	retrievedStruct, ok := retrievedValue.(TestStruct)
	assert.True(t, ok, "Expected successful type assertion")

	assert.Equal(t, value.ID, retrievedStruct.ID)
	assert.Equal(t, value.Name, retrievedStruct.Name)
	assert.Equal(t, value.Tags, retrievedStruct.Tags)
}

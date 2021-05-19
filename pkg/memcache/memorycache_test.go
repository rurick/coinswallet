package memorycache

import (
	"testing"
	"time"
)

const (
	testKey      string = "cache:test"
	testKeyEmpty string = "cache:empty"
	testValue    string = "Hello world"
)

// Test_Get get cache by key
func Test_Get(t *testing.T) {
	cache := New(10*time.Minute, 1*time.Hour)
	cache.Set(testKey, testValue, 1*time.Minute)
	value, found := cache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "Set and Get not simple:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	value, found = cache.Get(testKeyEmpty)
	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

// Test_Delete delete cache by key
func Test_Delete(t *testing.T) {
	cache := New(10*time.Minute, 1*time.Hour)
	cache.Set(testKey, testValue, 1*time.Minute)
	err := cache.Delete(testKey)

	if err != nil {
		t.Error("Error: ", "Cache delete failed")
	}

	value, found := cache.Get(testKey)
	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

	err = cache.Delete(testKeyEmpty)
	if err == nil {
		t.Error("Error: ", "An empty cache should return an error")
	}

}

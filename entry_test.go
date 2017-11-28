package redisproxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestCacheEntry_Expired  Checks that a newly created entry is fresh, waits a few seconds for it to expire, and checks that it's expired
func TestCacheEntry_Expired(t *testing.T) {
	entry := testEntry()

	assert.True(t, entry.Fresh(), "Test entry is still good")

	time.Sleep(time.Second * 3)

	assert.False(t, entry.Fresh(), "Test entry is expired")
}

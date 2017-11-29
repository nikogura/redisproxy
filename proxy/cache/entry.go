package cache

import (
	"time"
)

// CacheEntry  a struct representing a single cached entry
type CacheEntry struct {
	Expires time.Time
	Value   interface{}
	Key     string
}

// Fresh  Returns true if the current time is less than the entry's expiration.  Returns false otherwise.
func (e *CacheEntry) Fresh() bool {
	ttl := time.Now().Sub(e.Expires)

	if ttl < 0 {
		return true
	}

	return false
}

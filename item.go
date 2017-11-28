package redisproxy

import "time"

// CacheEntry  a struct representing a single cached entry
type CacheEntry struct {
	expires time.Time
	value interface{}

}

// Fresh  Returns true if the current time is less than the entry's expiration.  Returns false otherwise.
func (e *CacheEntry) Fresh() (bool) {
	ttl := time.Now().Sub(e.expires)

	if ttl < 0 {
		return true
	}

	return false
}



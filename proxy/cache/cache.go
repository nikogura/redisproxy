package cache

import (
	"container/list"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
	"time"
)

// Cache The actual cache object
type Cache struct {
	sync.RWMutex
	Ttl             time.Duration
	Entries         map[string]*list.Element
	MaxEntries      int
	AgeList         *list.List
	FetchFunc       FetchFunc
	fetchLock       sync.Mutex
	FetchInProgress map[string]time.Time
	FetchTimeout    time.Duration
	logger          *log.Logger
	RedisAddr       string
}

// FetchFunc Fetcher function.  Implemented separately so that I can make a mock one for testing
type FetchFunc func(key string, redisAddr string) (value interface{}, err error)

// NewCache  Creates a new cache.  Requires arguments for maxEntries (number of items in the cache) and maxAge(How long something will reside in the cache)
func NewCache(maxEntries int, maxAge time.Duration, fetchFunc FetchFunc, fetchTimeout time.Duration, redisAddr string) *Cache {

	logger := log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	c := &Cache{
		Ttl:             maxAge,
		Entries:         make(map[string]*list.Element),
		MaxEntries:      maxEntries,
		FetchInProgress: make(map[string]time.Time),
		AgeList:         list.New(),
		FetchFunc:       fetchFunc,
		FetchTimeout:    fetchTimeout,
		logger:          logger,
		RedisAddr:       redisAddr,
	}

	return c
}

// Get Gets an item from the cache, or if it's not in the cache, tries to get it from redis.  Automatically removes oldest entries from entry list if we exceed the maxEntries limit.
func (c *Cache) Get(key string) (entry *CacheEntry, err error) {

	// lock so it doesn't get written while we're reading it

	c.RLock()
	element, exists := c.Entries[key]
	c.RUnlock()

	//  If it isn't in the cache, go get it.
	if !exists {
		c.logger.Printf("Item not in cache.  Fetching.")
		return c.Fetch(key)
	}

	c.logger.Printf("Retrieving item from cache.")

	// this will, of course blow chunks if the entry's value is not a CacheElement.
	entry, ok := element.Value.(*CacheEntry)
	if !ok {
		err = errors.New("Couldn't extract a CacheEntry from the list element.  Wtf did you put in there?")
		return entry, err
	}

	// If it *is* in the cache, return it if it's fresh, moving it to the head of the age list, since it's now the freshest.
	if entry.Fresh() {
		c.RLock()
		c.AgeList.MoveToFront(element)
		c.RUnlock()

		return entry, err
	}

	// At this point it is in the cache, but it's stale.  Get rid of it.
	// locking and unlocking are performed in the function.
	c.RemoveElement(element)

	// Get a fresh version
	entry, err = c.Fetch(key)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to fetch key %q", key))
		return entry, err
	}

	return entry, err
}

// RemoveElement Removes the element from the list and age list.
func (c *Cache) RemoveElement(element *list.Element) {
	entry, ok := element.Value.(*CacheEntry)
	if ok {
		c.logger.Printf("Purging %s from cache", entry.Key)
		c.RLock()
		key := entry.Key

		if _, ok := c.Entries[key]; ok {
			delete(c.Entries, key)
		}
		c.AgeList.Remove(element)

		c.RUnlock()
		return
	}

	c.logger.Printf("Failed to remove element.  WTF did you put in here?")

}

// Fetch What actually reaches out and gets stuff by locking the cache and running the fetch func
func (c *Cache) Fetch(key string) (entry *CacheEntry, err error) {
	now := time.Now()
	// fetch item
	c.fetchLock.Lock()

	// Are we already fetching it?
	start, exists := c.FetchInProgress[key]

	// if so, have we timed out?
	if exists && start.Add(c.FetchTimeout).Before(now) {
		c.logger.Printf("We are already fetching %s", key)
		// if so, screw it, return an error
		c.fetchLock.Unlock()
		err = errors.New(fmt.Sprintf("Timeout fetching %s.  is fetchTimeout too short?", key))

		return entry, err
	}

	c.FetchInProgress[key] = now
	c.fetchLock.Unlock()

	// actually get the thing we're looking for
	value, err := c.FetchFunc(key, c.RedisAddr)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to fetch %s", key))
		return entry, err
	}

	if value != nil { // dont' bother storing nil values.
		entry = &CacheEntry{
			Expires: now.Add(c.Ttl),
			Value:   value,
			Key:     key,
		}

		c.RLock()

		element := c.AgeList.PushFront(entry)

		c.Entries[key] = element

		// Finally, check to see if we're over the configured cache size
		if len(c.Entries) > c.MaxEntries {
			c.logger.Printf("Max entries of %d reached.", c.MaxEntries)
			c.logger.Printf("Too many entries.  Purging the eldest.")
			eldest := c.AgeList.Back()
			c.RemoveElement(eldest)
		}

		c.RUnlock()
	}

	c.fetchLock.Lock()
	delete(c.FetchInProgress, key)
	c.fetchLock.Unlock()

	return entry, err

}

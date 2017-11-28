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
	ttl             time.Duration
	entries         map[string]*list.Element
	maxEntries      int
	ageList         *list.List
	fetchFunc       FetchFunc
	fetchLock       sync.Mutex
	fetchInProgress map[string]time.Time
	fetchTimeout    time.Duration
	logger          *log.Logger
}

// FetchFunc Fetcher function.  Implemented separately so that I can make a mock one for testing
type FetchFunc func(key string) (value interface{}, err error)

// NewCache  Creates a new cache.  Requires arguments for maxEntries (number of items in the cache) and maxAge(How long something will reside in the cache)
func NewCache(maxEntries int, maxAge time.Duration, fetchFunc FetchFunc, fetchTimeout time.Duration) *Cache {

	logger := log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	c := &Cache{
		ttl:             maxAge,
		entries:         make(map[string]*list.Element),
		maxEntries:      maxEntries,
		fetchInProgress: make(map[string]time.Time),
		ageList:         list.New(),
		fetchFunc:       fetchFunc,
		fetchTimeout:    fetchTimeout,
		logger:          logger,
	}

	return c
}

// Get Gets an item from the cache, or if it's not in the cache, tries to get it from redis.  Automatically removes oldest entries from entry list if we exceed the maxEntries limit.
func (c *Cache) Get(key string) (entry *CacheEntry, err error) {

	// lock so it doesn't get written while we're reading it

	c.RLock()
	element, exists := c.entries[key]
	c.RUnlock()

	//  If it isn't in the cache, go get it.
	if !exists {
		c.logger.Printf("Item not in cache.  Fetching.")
		return c.Fetch(key)
	}

	// this will, of course blow chunks if the entry's value is not a CacheElement.
	entry, ok := element.Value.(*CacheEntry)
	if !ok {
		err = errors.New("Couldn't extract a CacheEntry from the list element.  Wtf did you put in there?")
		return entry, err
	}

	// If it *is* in the cache, return it if it's fresh, moving it to the head of the age list, since it's now the freshest.
	if entry.Fresh() {
		c.RLock()
		c.ageList.MoveToFront(element)
		c.RUnlock()

		return entry, err
	}

	// At this point it is in the cache, but it's stale.  Get rid of it.
	// locking and unlocking are performed in the function.
	c.RemoveElement(key, element)

	// Get a fresh version
	entry, err = c.Fetch(key)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to fetch key %q", key))
		return entry, err
	}

	return entry, err
}

// RemoveElement Removes the element and key from the list and age list.
func (c *Cache) RemoveElement(key string, element *list.Element) {
	c.RLock()
	if _, ok := c.entries[key]; ok {
		delete(c.entries, key)
	}
	c.ageList.Remove(element)
	c.RUnlock()
}

// Fetch What actually reaches out and gets stuff by locking the cache and running the fetch func
func (c *Cache) Fetch(key string) (entry *CacheEntry, err error) {
	c.logger.Printf("In Fetch\n")
	now := time.Now()
	// fetch item
	c.fetchLock.Lock()

	// Are we already fetching it?
	start, exists := c.fetchInProgress[key]

	// if so, have we timed out?
	if exists && start.Add(c.fetchTimeout).Before(now) {
		c.logger.Printf("We are already fetching %s", key)
		// if so, screw it, return an error
		c.fetchLock.Unlock()
		err = errors.New(fmt.Sprintf("Timeout fetching %s.  is fetchTimeout too short?", key))

		return entry, err
	}

	c.fetchInProgress[key] = now
	c.fetchLock.Unlock()

	c.logger.Printf("Calling fetchFunc\n")
	// actually get the thing we're looking for
	value, err := c.fetchFunc(key)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to fetch %s", key))
		return entry, err
	}

	if value != nil { // dont' bother storing nil values.
		entry = &CacheEntry{
			expires: now.Add(c.ttl),
			value:   value,
		}

		c.RLock()

		element := c.ageList.PushFront(entry)

		c.entries[key] = element

		// Finally, check to see if we're over the configured cache size
		c.logger.Printf("Max entries: %d", c.maxEntries)
		if len(c.entries) > c.maxEntries {
			c.logger.Printf("Too many entries.  Purging one.")
			// remove the oldest entry from the ageList
			c.RemoveElement(key, element)
		}

		c.RUnlock()
	}

	c.fetchLock.Lock()
	delete(c.fetchInProgress, key)
	c.fetchLock.Unlock()

	return entry, err

}

package cache

import "time"

// I like putting my test fixtures in their own file, and making functions to return even simple things like strings.

// This protects me from typing the strings over and over and potentially fat fingering them all over the place.

// Those kinds of errors are just annoying, and I'd rather spend my time elsewhere

// testInterval  The ttl for entries in our tests
func testInterval() time.Duration {
	return time.Second * 2
}

// testValue  The test value we'll assert against
func testValue() string {
	return "goongala"
}

// testEntry  A fully formed CacheEntry to be used for testing
func testEntry() CacheEntry {
	return CacheEntry{
		expires: time.Now().Add(testInterval()),
		value:   testValue(),
	}
}

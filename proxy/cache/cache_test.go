package cache

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	c := NewCache(3, time.Second*3, unitTestFetchFunc, time.Second*1, "")

	key := testFoo()

	expected := testFoo()

	entry, err := c.Get(key)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key, err)
		t.Fail()
	}

	actual := entry.Value

	assert.Equal(t, expected, actual, "fetched string matches expectations")

	ttl1 := time.Now().Sub(entry.Expires)

	time.Sleep(time.Second * 1)

	entry, err = c.Get(key)

	actual = entry.Value

	assert.Equal(t, expected, actual, "fetched string matches expectations")

	ttl2 := time.Now().Sub(entry.Expires)

	assert.True(t, ttl2 > ttl1, "Time to live is indeed winding down.")

	time.Sleep(time.Second * 2)

	assert.False(t, entry.Fresh(), "Entry has expired")

	entry, err = c.Get(key)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key, err)
		t.Fail()
	}

	assert.True(t, entry.Fresh(), "Newly fetched entry is fresh.")

	expectedNumber := 10
	entry, err = c.Get("ten")
	if err != nil {
		log.Printf("Error getting key: %s", err)
		t.Fail()
	}

	actualNumber := entry.Value

	assert.Equal(t, expectedNumber, actualNumber, "Fetchng numbers works too")

}

func TestCache_CacheLimit(t *testing.T) {
	c := NewCache(3, time.Second*3, unitTestFetchFunc, time.Second*1, "")

	key1 := testFoo()
	_, err := c.Get(key1)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key1, err)
		t.Fail()
	}

	assert.True(t, len(c.Entries) == 1, "one entry in cache")

	key2 := testBar()
	_, err = c.Get(key2)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key2, err)
		t.Fail()
	}

	log.Printf("%d entries in cache", len(c.Entries))

	assert.True(t, len(c.Entries) == 2, "two entries in cache")

	key3 := testWip()
	_, err = c.Get(key3)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key3, err)
		t.Fail()
	}

	log.Printf("%d entries in cache", len(c.Entries))

	assert.True(t, len(c.Entries) == 3, "three entries in cache")

	key4 := testZoz()
	_, err = c.Get(key4)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key4, err)
		t.Fail()
	}

	log.Printf("%d entries in cache", len(c.Entries))

	assert.True(t, len(c.Entries) == 3, "three entries in cache")
}

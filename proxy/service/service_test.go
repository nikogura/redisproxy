package service

import (
	"fmt"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var proxy *Proxy

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setUp() {
	log.Printf("Attempting to find a free port on which to run the service.\n")
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("Failed to get a free port: %s\n", err)
		os.Exit(1)
	}

	log.Printf("Running with port %d\n", port)
	log.Printf("Creating proxy with the following:\n\tPort: %d\n\tCapacity: %d\n\tAge: %d\r\tTimeout: %d\r\tRedis: %s\n", port, testCapacity(), testMaxAge(), testTimeout(), testRedisAddr())

	proxy = TestProxy(port, testCapacity(), testMaxAge(), testTimeout(), testRedisAddr(), integTestFetchFunc)

	log.Printf("Running proxy\n")

	go proxy.Run()

	log.Printf("Proxy is running.  Sleeping 5 seconds to let it get it's bearings before we hit it.\n")

	time.Sleep(time.Second * 5)

	log.Printf("Setup complete.\n")

}

func tearDown() {

}

// This is basically the same test as in cache_test, but this one is spinning up the actual http service and exercising it.
func TestCacheAndLimit(t *testing.T) {
	// ************** Entry One **************************
	key1 := testFoo()

	uri := fmt.Sprintf("http://localhost%s/%s", proxy.Port, key1)
	log.Printf("Getting %s", uri)

	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key1, err)
		t.Fail()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		t.Fail()
	}

	assert.Equal(t, fmt.Sprintf("\"%s\"\n", key1), string(body), "Http response for key meets expectations.")

	assert.True(t, len(proxy.Cache.Entries) == 1, "one entry in cache")

	// ************** Entry Two **************************

	key2 := testBar()

	uri = fmt.Sprintf("http://localhost%s/%s", proxy.Port, key2)
	log.Printf("Getting %s", uri)

	resp, err = http.Get(uri)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key2, err)
		t.Fail()
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		t.Fail()
	}

	assert.Equal(t, fmt.Sprintf("\"%s\"\n", key2), string(body), "Http response for key meets expectations.")

	assert.True(t, len(proxy.Cache.Entries) == 2, "two entries in cache")

	// ************** Entry Three **************************
	key3 := testWip()

	uri = fmt.Sprintf("http://localhost%s/%s", proxy.Port, key3)
	log.Printf("Getting %s", uri)

	resp, err = http.Get(uri)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key3, err)
		t.Fail()
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		t.Fail()
	}

	assert.Equal(t, fmt.Sprintf("\"%s\"\n", key3), string(body), "Http response for key meets expectations.")

	assert.True(t, len(proxy.Cache.Entries) == 3, "three entries in cache")

	// ************** Entry Four **************************
	key4 := testZoz()

	uri = fmt.Sprintf("http://localhost%s/%s", proxy.Port, key4)
	log.Printf("Getting %s", uri)

	resp, err = http.Get(uri)
	if err != nil {
		log.Printf("Error fetching key %s: %s", key4, err)
		t.Fail()
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		t.Fail()
	}

	assert.Equal(t, fmt.Sprintf("\"%s\"\n", key4), string(body), "Http response for key meets expectations.")

	assert.True(t, len(proxy.Cache.Entries) == 3, "three entries in cache")

}

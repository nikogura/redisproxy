package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Proxy struct to represent the proxy server itself
type Proxy struct {
	Cache     *Cache
	RedisAddr string
	Port      string
}

// NewProxy creates, guess what?  a new proxy
func NewProxy(port int, maxEntries int, maxAge int, timeout int, redisAddr string) *Proxy {
	portString := strconv.Itoa(port)

	realPort := fmt.Sprintf(":%s", portString)

	proxy := &Proxy{
		Cache:     NewCache(maxEntries, time.Duration(maxAge)*time.Second, Fetcher, time.Duration(timeout)*time.Second, redisAddr),
		Port:      realPort,
		RedisAddr: redisAddr,
	}

	return proxy
}

// Run actually runs the http server for the proxy.  It does not detatch from the console
func (p *Proxy) Run() (err error) {
	http.HandleFunc("/", p.Handle)
	err = http.ListenAndServe(p.Port, nil)

	return err
}

// Handle is the http handler for all incoming requests.
func (p *Proxy) Handle(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.RequestURI, "/")

	log.Printf("Received request for %s\n", key)

	entry, err := p.Cache.Get(key)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	log.Printf("Get result: %s", entry.Value)

	if entry != nil {
		value := entry.Value

		valtype := reflect.TypeOf(value).String()

		if valtype == "string" {
			fmt.Fprintf(w, "%q\n", value)
		} else {
			fmt.Fprintf(w, "(%s) %s\n", valtype, value)
		}
		log.Printf("Done with request\n")

		return
	}

	log.Printf("No result, sending nil\n")

	fmt.Fprint(w, "(nil)\n")
	log.Printf("Done with request\n")

}

// Fetcher The function that actually gets info from redis.  This is used when the proxy is run for reals.  In testing it's replaced by an in memory function reading from a test fixture
func Fetcher(key string, redisAddr string) (value interface{}, err error) {
	r := regexp.MustCompile(`.+:\d+`)

	var fqRedisAddr string

	if r.MatchString(redisAddr) {
		fqRedisAddr = redisAddr
	} else {
		fqRedisAddr = fmt.Sprintf("%s:6379", redisAddr)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fqRedisAddr,
		Password: "",
		DB:       0,
	})

	fetchedval, err := client.Get(key).Result()
	if err == redis.Nil {
		return value, nil
	} else if err != nil {
		return value, err
	}

	value = fetchedval

	return value, err
}

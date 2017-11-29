// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nikogura/redisproxy/cache"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the proxy with the supplied options",
	Long: `
Runs the proxy with the supplied options.

Does not detatch from the console.
`,
	Run: func(cmd *cobra.Command, args []string) {
		portString := strconv.Itoa(cachePort)

		port := fmt.Sprintf(":%s", portString)

		log.Printf("Starting Cache on port %s\n", port)
		log.Printf("Cache Expiration: %d seconds\n", cacheExpirationSeconds)
		log.Printf("Cache Capacity: %d entries\n", cacheCapacity)
		log.Printf("Upstream Redis Instance: %q\n", redisAddr)

		c := cache.NewCache(cacheCapacity, time.Duration(cacheExpirationSeconds)*time.Second, fetcher, time.Second*5)

		http.Handle("/", Handler{cache: c})
		http.ListenAndServe(port, nil)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

}

type Handler struct {
	cache *cache.Cache
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Howdy\n")

	entry, err := h.cache.Get("foo")
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
	}

	if entry != nil {
		value := entry.Value
		fmt.Fprint(w, value)
	}

	fmt.Fprint(w, "nil")
}

func fetcher(key string) (value interface{}, err error) {
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
		return value, err
	} else if err != nil {
		return value, err
	}

	value = fetchedval

	return value, err
}

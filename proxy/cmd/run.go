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
	"github.com/nikogura/redisproxy/proxy/service"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

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

		proxy := service.NewProxy(cachePort, cacheCapacity, cacheExpirationSeconds, 5, redisAddr)

		err := proxy.Run()
		if err != nil {
			log.Fatalf("Error Running proxy: %s", err)
		}

	},
}

func init() {
	RootCmd.AddCommand(runCmd)

}

//type Proxy struct {
//	cache *cache.Cache
//}
//
//func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
//	key := strings.TrimPrefix(r.RequestURI, "/")
//
//	log.Printf("Received request for %s\n", key)
//
//	entry, err := p.cache.Get(key)
//	if err != nil {
//		fmt.Fprintf(w, "Error: %s\n", err)
//		return
//	}
//
//	log.Printf("Get result: %s", entry.Value)
//
//	if entry != nil {
//		value := entry.Value
//
//		valtype := reflect.TypeOf(value).String()
//
//		if valtype == "string" {
//			fmt.Fprintf(w, "%q\n", value)
//		} else {
//			fmt.Fprintf(w, "(%s) %s\n", valtype, value)
//		}
//		log.Printf("Done with request\n")
//
//		return
//	}
//
//	log.Printf("No result, sending nil\n")
//
//	fmt.Fprint(w, "(nil)\n")
//	log.Printf("Done with request\n")
//
//}
//
//// Fetcher The function that actually gets info from redis.
//func Fetcher(key string) (value interface{}, err error) {
//	r := regexp.MustCompile(`.+:\d+`)
//
//	var fqRedisAddr string
//
//	if r.MatchString(redisAddr) {
//		fqRedisAddr = redisAddr
//	} else {
//		fqRedisAddr = fmt.Sprintf("%s:6379", redisAddr)
//	}
//
//	client := redis.NewClient(&redis.Options{
//		Addr:     fqRedisAddr,
//		Password: "",
//		DB:       0,
//	})
//
//	fetchedval, err := client.Get(key).Result()
//	if err == redis.Nil {
//		return value, nil
//	} else if err != nil {
//		return value, err
//	}
//
//	value = fetchedval
//
//	return value, err
//}

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

	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strconv"
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
		if redisAddr == "" {
			log.Fatalf("Cannot run without an IP address for a Redis instance.  Try again with: -r <redis ip:port>")
		}

		portString := strconv.Itoa(cachePort)

		port := fmt.Sprintf(":%s", portString)

		log.Printf("Starting Cache on port %s\n", port)
		log.Printf("Cache Expiration: %d seconds\n", cacheExpirationSeconds)
		log.Printf("Cache Capacity: %d entries\n", cacheCapacity)
		log.Printf("Upstream Redis IP Address: %s\n", redisAddr)

		http.HandleFunc("/", handler)
		http.ListenAndServe(port, nil)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Howdy\n")
}
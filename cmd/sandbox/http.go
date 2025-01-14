/*
Copyright © 2025 ks6088ts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package sandbox

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run HTTP Server",
	Long:  `This is a sandbox command to run a simple HTTP server.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags
		port, err := cmd.Flags().GetInt("port")
		// handle error
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/healthz" {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		})

		http.Handle("/metrics", promhttp.Handler())

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, world!")
		})

		http.HandleFunc("/flaky", func(w http.ResponseWriter, r *http.Request) {
			// Ensure POST method
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintln(w, "Method Not Allowed")
				return
			}

			// Define request structure
			type flakyRequest struct {
				Percent int `json:"percent"`
			}

			// Decode JSON
			var reqBody flakyRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid JSON format in request body")
				return
			}

			// Use default if not specified
			if reqBody.Percent <= 0 || reqBody.Percent > 100 {
				reqBody.Percent = 50
			}

			// Return 500 if random < reqBody.Percent
			if rand.Intn(100) < reqBody.Percent {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Internal Server Error")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Success")
		})

		log.Printf("Starting server on port %d\n", port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	sandboxCmd.AddCommand(httpCmd)

	httpCmd.Flags().IntP("port", "p", 8080, "Port number to listen")
}

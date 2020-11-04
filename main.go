// Copyright © 2018 Banzai Cloud
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

func main() {
	mux := http.NewServeMux()

	{
		exporter, err := prometheus.NewExporter(prometheus.Options{})
		if err != nil {
			log.Fatal(err)
		}

		view.RegisterExporter(exporter)
		mux.Handle("/metrics", exporter)
	}

	// Register stat views
	err := view.Register(
		// HTTP
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	mux.Handle("/demo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, "hello world")

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(fmt.Sprintf(`Welcome to Banzai Cloud Pipeline!

Your secret is: %s
`, os.Getenv("SECRET"))))
	}))

	server := &http.Server{
		Addr: ":8080",
		Handler: &ochttp.Handler{
			Handler: mux,
		},
	}

	log.Println("starting application")

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

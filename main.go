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

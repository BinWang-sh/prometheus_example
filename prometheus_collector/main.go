package main

import (
	apiCollector "binTest/prometheusTest/prometheus_collector/collector"
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {
	flag.Parse()

	apiC := apiCollector.NewApiCollector("demo")
	registry := prometheus.NewRegistry()
	registry.MustRegister(apiC)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	//http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>A Prometheus Exporter</title></head>
				<body>
				<h1>A Prometheus Exporter</h1>
				<p><a href='/metrics'>Metrics</a></p>
				</body>
				</html>`))
	})

	log.Printf("Starting Server at http://localhost:%s%s", *addr, "/metrics")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

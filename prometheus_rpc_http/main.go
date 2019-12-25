package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

var (
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "rpc_durations_seconds",
			Help: "RPC latency distributions.",
			//Objectives: map[float64]float64{0.5: 0.5, 0.9: 1.0, 0.99: 1.5},
		},
		[]string{"service"},
	)

	httpReqDurationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_req_durations_histogram",
			Help: "http req latency distributions.",
			// 4 buckets, starting from 0.1 and adding 0.5 between each bucket
			Buckets: prometheus.LinearBuckets(0.1, 0.5, 4),
		},
		[]string{"http_req_histogram"},
	)

	rpcCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_counter",
			Help: "RPC counts",
		},
		[]string{"api"},
	)

	rpcReqSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rpc_req_size",
			Help: "RPC request size",
		},
		[]string{"api"},
	)
)

func init() {
	// Register the summary, histogram, counter, gauge with Prometheus's default registry.
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(httpReqDurationsHistogram)
	prometheus.MustRegister(rpcCounter)
	prometheus.MustRegister(rpcReqSize)

}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			v := rand.Float64()
			rpcDurations.WithLabelValues("user_rpc").Observe(v)
			httpReqDurationsHistogram.WithLabelValues("booksvc_req").Observe(1.5 * v)

			rpcCounter.WithLabelValues("api_bookcontent").Add(float64(rand.Int31n(50)))
			rpcCounter.WithLabelValues("api_chapterlist").Add(float64(rand.Int31n(10)))

			rpcReqSize.WithLabelValues("api_bookcontent").Set(float64(rand.Int31n(8000)))
			rpcReqSize.WithLabelValues("api_chapterlist").Set(float64(rand.Int31n(5000)))
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			v := 0.5 + rand.Float64()
			rpcDurations.WithLabelValues("book_rpc").Observe(v)
			httpReqDurationsHistogram.WithLabelValues("booksvc_req").Observe(1.5 * v)

			rpcCounter.WithLabelValues("api_bookcontent").Add(float64(rand.Int31n(10)))
			rpcCounter.WithLabelValues("api_chapterlist").Add(float64(rand.Int31n(20)))

			time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
		}
	}()

	go func() {
		for {
			v := 1.0 + rand.Float64()
			rpcDurations.WithLabelValues("bookshelf_rpc").Observe(v)
			httpReqDurationsHistogram.WithLabelValues("booksvc_req").Observe(1.5 * v)

			rpcCounter.WithLabelValues("api_chapterlist").Add(float64(rand.Int31n(250)))
			rpcCounter.WithLabelValues("api_bookcontent").Add(float64(rand.Int31n(350)))

			time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

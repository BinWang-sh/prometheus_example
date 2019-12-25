package collector

import (
	"math/rand"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MatricType int32

const (
	MatricType_Counter   MatricType = 0
	MatricType_Gauge     MatricType = 1
	MatricType_Histogram MatricType = 2
	MatricType_Summary   MatricType = 3
)

func (p MatricType) String() string {
	switch p {
	case MatricType_Counter:
		return "Counter"
	case MatricType_Gauge:
		return "Gauge"
	case MatricType_Histogram:
		return "Histogram"
	case MatricType_Summary:
		return "Summary"
	default:
		return "UNKNOWN"
	}
}

type ApiCollector struct {
	nameSpace string
	mMetrics  map[string]*prometheus.Desc
	mutex     sync.Mutex
}

func NewApiCollector(newNamespace string) *ApiCollector {

	return &ApiCollector{
		nameSpace: newNamespace,
		mMetrics: map[string]*prometheus.Desc{
			"req_counter_metric":    prometheus.NewDesc("req_counter_metric", "The request counter matric", []string{"api"}, nil),
			"req_time_gauge_metric": prometheus.NewDesc("req_time_gauge_metric", "The request cost time matric", []string{"api"}, nil),
		},
	}
}

func (c *ApiCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.mMetrics {
		ch <- m
	}
}

func (c *ApiCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mockCounterMetricData := c.GenerateData(MatricType_Counter)
	for host, value := range mockCounterMetricData {
		ch <- prometheus.MustNewConstMetric(c.mMetrics["req_counter_metric"], prometheus.CounterValue, float64(value), host)
	}

	mockGaugeMetricData := c.GenerateData(MatricType_Gauge)
	for host, currentValue := range mockGaugeMetricData {
		ch <- prometheus.MustNewConstMetric(c.mMetrics["req_time_gauge_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}

}

func (c *ApiCollector) GenerateData(mtype MatricType) map[string]int {
	switch mtype.String() {
	case "Counter":
		mockCounterMetricData := map[string]int{
			"api/bookcontent": int(rand.Int31n(500000)),
			"api/chapterlist": int(rand.Int31n(50000)),
			"api/bookstore":   int(rand.Int31n(800000)),
		}
		return mockCounterMetricData
	case "Gauge":
		mockGaugeMetricData := map[string]int{
			"api/bookcontent": int(rand.Int31n(200)),
			"api/chapterlist": int(rand.Int31n(200)),
			"api/bookstore":   int(rand.Int31n(200)),
		}
		return mockGaugeMetricData
	}

	return nil
}

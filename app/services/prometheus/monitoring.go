package prometheus

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var reqCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "request_path",
	Subsystem: "middleware",
	Name:      "total_requests",
	Help:      "Total Requests recieved by this machine",
}, []string{"host", "path"})

var resCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "response_code",
	Subsystem: "middleware",
	Name:      "response_codes",
	Help:      "Count of various response code on the system",
},
	[]string{"code", "length"})

//IncrementRequestCount - increments the number of requests attended to by the system per path
func IncrementRequestCount(host, path string) {
	reqCount.WithLabelValues(host, path).Inc()
}

//IncrementResponseCount - Keeps a metrics of the responses returned by the system
func IncrementResponseCount(code, length int) {
	resCount.WithLabelValues(strconv.Itoa(code), strconv.Itoa(length)).Inc()
}

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vchitai/gomw"
)

var httpLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http.server.duration",
	Buckets: prometheus.DefBuckets, // seconds
}, []string{"operation"})

type ctxKey string

type latencyMeasure struct {
	startTime time.Time
}

func LatencyRecordingHTTPMiddleware(operation string) gomw.HTTPMiddleware {
	return gomw.NewHTTPMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		measure := latencyMeasure{startTime: time.Now()}
		newCtx := context.WithValue(request.Context(), ctxKey("latency-measure"), &measure)
		request = request.WithContext(newCtx)
		return request, true
	}, func(response gomw.HTTPResponse, request *http.Request) gomw.HTTPResponse {
		measure, ok := request.Context().Value(ctxKey("latency-measure")).(*latencyMeasure)
		if !ok {
			return response
		}
		hist, err := httpLatencyHistogram.GetMetricWithLabelValues(operation)
		if err != nil {
			return response
		}
		hist.Observe(time.Since(measure.startTime).Seconds())
		return response
	})
}

func helloWorld() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello world!"))
	}
}

func main() {
	prometheus.MustRegister(httpLatencyHistogram)
	var mux = http.NewServeMux() // create new server
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", LatencyRecordingHTTPMiddleware("hello world")(helloWorld()))

	l, err := net.Listen("tcp", ":10080")
	if err != nil {
		log.Fatal(err)
	}

	if err := http.Serve(l, mux); err != nil {
		log.Fatal(err)
	}
}

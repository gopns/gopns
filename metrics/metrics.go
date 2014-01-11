package metrics

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"time"
)

func GetCallMeters(metricName string) (callMeter metrics.Meter, errorMeter metrics.Meter) {
	callMeter = metrics.GetOrRegisterMeter(metricName, metrics.DefaultRegistry)
	errorMeter = metrics.GetOrRegisterMeter(metricName+"error", metrics.DefaultRegistry)

	return callMeter, errorMeter
}

func StartMetricPrinter() {
	go metrics.Log(metrics.DefaultRegistry, time.Minute, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
}

package metrics

import (
	gometrics "github.com/rcrowley/go-metrics"
	"testing"
	"time"
)

func TestGetCallMeters(t *testing.T) {

	callMeter, errorMeter := GetCallMeters("Test")
	callMeter.Mark(5)
	callMeter.Mark(2)

	errorMeter.Mark(3)
	if callMeter.Count() != 7 {
		t.Errorf("Meter Count does not line up for call meter, expected 7 got %d", callMeter.Count())
	}

	callMeterCompare := gometrics.DefaultRegistry.Get("Test")

	if callMeterCompare.(gometrics.Meter).Count() != 7 {
		t.Errorf("Meter Count does not line up for call meter, expected 7 got %d", callMeter.Count())
	}

	errorMeterCompare := gometrics.DefaultRegistry.Get("Test.error")
	errorMeterCompare.(gometrics.Meter).Mark(2)

	if errorMeterCompare.(gometrics.Meter).Count() != 5 {
		t.Errorf("Error Meter Count does not line up for call meter, expected 5 got %d", callMeter.Count())
	}

}

func TestTimer(t *testing.T) {

	timer := GetTimer("TestTimer")
	timer.Time(func() {
		time.Sleep(time.Second)
	})

	if timer.Count() != 1 {
		t.Errorf("Error Timer Count does not line up, expected 1 got %d", timer.Count())
	}

	if timer.Mean() < time.Second.Seconds() {
		t.Errorf("Error Timer mean does not line up, expected > 1 got %f", timer.Mean())
	}

}

func TestReporters(t *testing.T) {

	StartGraphiteReporter("http://www.google.com", "key", "prefix")
	StartMetricPrinter()
}

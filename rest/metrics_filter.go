package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/metrics"
	"regexp"
)

var dotChars = regexp.MustCompile(`[/\\]`)
var removeChars = regexp.MustCompile(`[\(\) ?<>;:~@#$%^&*+=]`)

type MetricsFilter interface {
	MetricsCollector()
}

func NewTimingFilter(key string) restful.FilterFunction {
	tFilter := new(timingMetricsFilter)

	tFilter.key = sanitizeKey(key)
	tFilter.timer = metrics.GetTimer(key)
	return tFilter.MetricsCollector
}

func sanitizeKey(key string) string {
	key = dotChars.ReplaceAllString(key, ".")
	return removeChars.ReplaceAllString(key, "")
}

type timingMetricsFilter struct {
	key   string
	timer metrics.Timer
}

func (this *timingMetricsFilter) MetricsCollector(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	this.timer.Time(func() {
		chain.ProcessFilter(req, resp)
	})

}

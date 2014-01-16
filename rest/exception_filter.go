package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/exception"
	"net/http"
	//"github.com/gopns/gopns/metrics"
)

func ExceptionFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer recoverPanic(resp)
	chain.ProcessFilter(req, resp)

}

type responseWriter interface {
	WriteErrorString(int, string) error
}

func recoverPanic(resp responseWriter) {
	recovered := recover()

	if recovered == nil {
		// Do Nothing
		return
	} else if webException, match := recovered.(exception.WebException); match {
		resp.WriteErrorString(webException.ResponseStatus(), webException.Message())
	} else if exception_, match := recovered.(exception.Exception); match {
		resp.WriteErrorString(http.StatusInternalServerError, exception_.Message())
	} else if message, match := recovered.(string); match {
		resp.WriteErrorString(http.StatusInternalServerError, message)
	}

}

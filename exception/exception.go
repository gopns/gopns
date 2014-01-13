package exception

import (
	"net/http"
)

type Exception interface {
	Message() string
}

type WebException interface {
	Exception
	ResponseStatus() int
}

type baseException struct {
	message string
}

type baseWebException struct {
	baseException
	responseStatus int
}

func (this *baseException) Message() string {
	return this.message
}

func (this *baseWebException) ResponseStatus() int {
	return this.responseStatus
}

func NewException(msg string) Exception {
	return &baseException{msg}
}

func BadRequestException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusBadRequest}
}

func UnauthorizedException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusUnauthorized}
}

func PaymentRequiredException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusPaymentRequired}
}

func ForbiddenException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusForbidden}
}

func NotFoundException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusNotFound}
}

func MethodNotAllowedException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusMethodNotAllowed}
}

func NotAcceptableException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusNotAcceptable}
}

func ProxyAuthRequiredException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusProxyAuthRequired}
}

func ConflictException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusConflict}
}

func GoneException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusGone}
}

func LengthRequiredException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusLengthRequired}
}

func PreconditionFailedException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusPreconditionFailed}
}

func RequestEntityTooLargeException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusRequestEntityTooLarge}
}

func RequestURITooLongException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusRequestURITooLong}
}

func UnsupportedMediaTypeException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusUnsupportedMediaType}
}

func RequestedRangeNotSatisfiableException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusRequestedRangeNotSatisfiable}
}

func ExpectationFailedException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusExpectationFailed}
}

func InternalServerErrorException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusInternalServerError}
}

func NotImplemented(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusNotImplemented}
}

func BadGatewayException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusBadGateway}
}

func ServiceUnavailableException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusServiceUnavailable}
}

func GatewayTimeoutException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusGatewayTimeout}
}

func HTTPVersionNotSupportedException(msg string) Exception {
	return &baseWebException{baseException{msg}, http.StatusHTTPVersionNotSupported}
}

func ConditionalThrowNewException(err error) {
	if err != nil {
		panic(Exception(&baseException{err.Error()}))
	}
}

func ConditionalThrowBadRequestException(err error) {
	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusBadRequest}))
	}
}

func ConditionalThrowUnauthorizedException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusUnauthorized})))
	}
}

func ConditionalThrowPaymentRequiredException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusPaymentRequired})))
	}
}

func ConditionalThrowForbiddenException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusForbidden})))
	}
}

func ConditionalThrowNotFoundException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusNotFound})))
	}
}

func ConditionalThrowMethodNotAllowedException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusMethodNotAllowed})))
	}
}

func ConditionalThrowNotAcceptableException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusNotAcceptable})))
	}
}

func ConditionalThrowProxyAuthRequiredException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusProxyAuthRequired})))
	}
}

func ConditionalThrowConflictException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusConflict})))
	}
}

func ConditionalThrowGoneException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusGone})))
	}
}

func ConditionalThrowLengthRequiredException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusLengthRequired})))
	}
}

func ConditionalThrowPreconditionFailedException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusPreconditionFailed})))
	}
}

func ConditionalThrowRequestEntityTooLargeException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusRequestEntityTooLarge})))
	}
}

func ConditionalThrowRequestURITooLongException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusRequestURITooLong})))
	}
}

func ConditionalThrowUnsupportedMediaTypeException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusUnsupportedMediaType})))
	}
}

func ConditionalThrowRequestedRangeNotSatisfiableException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusRequestedRangeNotSatisfiable})))
	}
}

func ConditionalThrowExpectationFailedException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusExpectationFailed})))
	}
}

func ConditionalThrowInternalServerErrorException(err error) {

	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusInternalServerError}))
	}
}

func ConditionalThrowNotImplemented(err error) {
	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusNotImplemented}))
	}
}

func ConditionalThrowBadGatewayException(err error) {
	if err != nil {
		panic(WebException(&(baseWebException{baseException{err.Error()}, http.StatusBadGateway})))
	}
}

func ConditionalThrowServiceUnavailableException(err error) {
	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusServiceUnavailable}))
	}
}

func ConditionalThrowGatewayTimeoutException(err error) {
	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusGatewayTimeout}))
	}
}

func ConditionalThrowHTTPVersionNotSupportedException(err error) {
	if err != nil {
		panic(WebException(&baseWebException{baseException{err.Error()}, http.StatusHTTPVersionNotSupported}))
	}
}

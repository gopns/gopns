package rest

import (
	"errors"
	"github.com/gopns/gopns/exception"
	"strings"
	"testing"
)

//Tests for Metrics Filter
func TestKeySanitization(t *testing.T) {

	sanitizedKey := sanitizeKey("()\\/?<>;:~`@#$%^&*+=")
	if strings.ContainsAny(sanitizedKey, " ,()\\/?<>;:~@#$%^&*+=") {
		t.Error("Invalid sanitization. Returned string :", sanitizedKey)
	}
}

//Tests for Exception Filter
type RespTestWriter struct {
	t               *testing.T
	expectedStatus  int
	expectedMessage string
}

func (this *RespTestWriter) WriteErrorString(status int, message string) error {
	if status != this.expectedStatus {
		this.t.Errorf("Inavlid stats %d, %d expected", status, this.expectedStatus)
	} else if message != this.expectedMessage {
		this.t.Errorf("Inavlid message %s, %s expected", message, this.expectedMessage)
	}
	return nil
}

func recoverPanicFail(t *testing.T) {
	if recover() != nil {
		//Did Not catch panic so fail test
		t.Error("Unable to recover from panic when needed")
	}
}

func TestRecoverPanicString(t *testing.T) {
	defer recoverPanicFail(t)
	sendPanicString(t)
}

func sendPanicString(t *testing.T) {
	defer recoverPanic(&RespTestWriter{t, 500, "Panic String"})
	panic("Panic String")
}

func TestRecoverPanicException(t *testing.T) {
	defer recoverPanicFail(t)
	sendPanicException(t)
}

func sendPanicException(t *testing.T) {
	defer recoverPanic(&RespTestWriter{t, 500, "Exception Message"})
	panic(exception.NewException("Exception Message"))
}

func TestRecoverPanicWebException(t *testing.T) {
	defer recoverPanicFail(t)
	sendPanicWebException(t)
}

func sendPanicWebException(t *testing.T) {
	defer recoverPanic(&RespTestWriter{t, 400, "Bad Requesst Message"})
	panic(exception.BadRequestException("Bad Requesst Message"))
}

func TestRecoverConditionalPanicTrue(t *testing.T) {
	defer recoverPanicFail(t)
	sendPanicConditionalTrue(t)
}

func sendPanicConditionalTrue(t *testing.T) {
	defer recoverPanic(&RespTestWriter{t, 401, "Unauthorized"})
	exception.ConditionalThrowUnauthorizedException(errors.New("Unauthorized"))
}

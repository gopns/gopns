package restutil

import (
	"code.google.com/p/gorest"
)

type RestError struct {
	ErrorCase       bool
	ErrorCode       int
	ErrorMessage    string
	ResponseBuilder *gorest.ResponseBuilder
}

func GetRestError(responseBuilder *gorest.ResponseBuilder) *RestError {
	return &RestError{false, 200, "", responseBuilder}
}

func CheckError(err error, restError *RestError, errorCode int) {
	if err != nil {
		restError.ErrorCase = true
		restError.ErrorCode = errorCode
		restError.ErrorMessage = err.Error()

		panic(err.Error())
	}
}

func HandleErrors(restError *RestError) {

	if restError.ErrorCase && recover() != nil {
		restError.ResponseBuilder.SetResponseCode(restError.ErrorCode).WriteAndOveride([]byte("{\"message\":\"" + restError.ErrorMessage + "\"}"))
	}
}

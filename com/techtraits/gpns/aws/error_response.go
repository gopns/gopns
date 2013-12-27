package aws

type ErrorResponse struct {
	Error ErrorStruct
}

type ErrorStruct struct {
	Type    string
	Code    string
	Message string
}

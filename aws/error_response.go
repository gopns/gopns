package aws

type ErrorResponse struct {
	Error ErrorStruct
}

type ErrorStruct struct {
	Type    string `json:"__type"`
	Code    string
	Message string `json:"message"`
}

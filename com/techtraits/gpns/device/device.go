package device

type Device struct {
	Alias    string
	Id       string
	Arn      string
	Platform string
	Locale   string
	Tags     []string
}

package sns

import (
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"net/http"
	"net/url"
	"strings"
)

func MakeRequest(host string, values url.Values, platformAppName string) (*http.Response, error) {
	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	aws.SignRequest(req, "sns", platformAppName)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

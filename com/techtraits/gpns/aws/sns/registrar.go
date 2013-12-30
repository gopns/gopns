package sns

import (
	"encoding/xml"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var registrarInstance RegistrarStruct

func RegistrarInstance() Registrar {
	return registrarInstance
}

type RegistrarStruct struct {
}

type Registrar interface {
	RegisterDevice(platformAppName string, token string, customData string) (arn string, err error)
}

func (this RegistrarStruct) RegisterDevice(platformAppName string, token string, customData string) (arn string, err error) {
	values := url.Values{}
	values.Set("Action", "CreatePlatformEndpoint")
	values.Set("CustomUserData", customData)
	values.Set("Token", token)
	values.Set("PlatformApplicationArn", gpnsconfig.AWSConfigInstance().PlatformApps()[platformAppName].Arn())
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	url_, err := url.Parse("http://sns." + gpnsconfig.AWSConfigInstance().PlatformApps()[platformAppName].Region() + ".amazonaws.com/")
	if err != nil {
		return "", err
	}

	aws.SignRequest("POST", "/", values, url_.Host)

	response, err := http.PostForm(url_.String(), values)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return "", errors.New("Unable to register device. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var createResponse CreateResponse
		xml.Unmarshal(content, &createResponse)
		return createResponse.CreatePlatformEndpointResult.EndpointArn, nil
	}
}

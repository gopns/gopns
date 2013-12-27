package sns

import (
	"encoding/xml"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RegistrarStruct struct {
	awsConfig gpnsconfig.AWSConfig
}

func Initilize(awsConfig gpnsconfig.AWSConfig) Registrar {
	return RegistrarStruct{awsConfig}
}

type Registrar interface {
	RegisterDevice(platformAppName string, token string, customData string) (arn string, err error)
}

func (this RegistrarStruct) RegisterDevice(platformAppName string, token string, customData string) (arn string, err error) {
	values := url.Values{}
	values.Set("Action", "CreatePlatformEndpoint")
	values.Set("CustomUserData", customData)
	values.Set("Token", token)
	values.Set("PlatformApplicationArn", this.awsConfig.PlatformApps()[platformAppName].Arn())
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	url_, err := url.Parse("http://sns." + this.awsConfig.PlatformApps()[platformAppName].Region() + ".amazonaws.com/")
	if err != nil {
		return "", err
	}

	aws.SignRequest(this.awsConfig, "POST", "/", values, url_.Host)

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
	return "", err
}

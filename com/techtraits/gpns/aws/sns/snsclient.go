package sns

import (
	"encoding/xml"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RegisterDevice(token string, customData string, userId string,
	userSecret string, region string, applicationArn string) (arn string, err error) {
	values := url.Values{}
	values.Set("Action", "CreatePlatformEndpoint")
	values.Set("CustomUserData", customData)
	values.Set("Token", token)
	values.Set("PlatformApplicationArn", applicationArn)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest("http://sns."+region+".amazonaws.com/",
		values, userId, userSecret, region)

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

func PublishNotification(arn string, title string, message string, userId string,
	userSecret string, region string, applicationType string) (err error) {

	values := url.Values{}
	values.Set("Action", "Publish")
	values.Set("Message", formatMessage(applicationType, title, message))
	values.Set("MessageStructure", "json")
	values.Set("TargetArn", arn)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest("http://sns."+region+".amazonaws.com/",
		values, userId, userSecret, region)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to send Push Notification. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		return nil
	}

}

func formatMessage(appType string, title string, message string) string {
	formattedMessage := "{\"" + appType + "\":\"{\\\"data\\\":{\\\"id\\\":\\\"MESSAGE_ID\\\",\\\"title\\\":\\\"" + title + "\\\",\\\"alert\\\":\\\"" + message + "\\\"}}\"}"
	return formattedMessage
}

func makeRequest(host string, values url.Values, userId string,
	userSecret string, region string) (*http.Response, error) {

	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	aws.SignRequest(req, userId, userSecret, "sns", region)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

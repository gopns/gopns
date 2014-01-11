package sns

import (
	"encoding/xml"
	"errors"
	"github.com/gopns/gopns/aws"
	"github.com/gopns/gopns/metrics"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SNSClient interface {
	RegisterDevice(token string, customData string, applicationArn string) (arn string, err error)
	PublishNotification(arn string, title string, message string, applicationType string) (err error)
}

type BasicSNSClient struct {
	UserId     string
	UserSecret string
	Region     string
}

func New(
	userId string,
	userSecret string,
	region string) (SNSClient, error) {
	return &BasicSNSClient{userId, userSecret, region}, nil

}

func (this *BasicSNSClient) RegisterDevice(token string, customData string, applicationArn string) (arn string, err error) {

	callMeter, errorMeter := metrics.GetCallMeters("sns.register_device")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "CreatePlatformEndpoint")
	values.Set("CustomUserData", customData)
	values.Set("Token", token)
	values.Set("PlatformApplicationArn", applicationArn)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest("http://sns."+this.Region+".amazonaws.com/",
		values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return "", errors.New("Unable to register device. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var createResponse CreateResponse
		xml.Unmarshal(content, &createResponse)
		return createResponse.CreatePlatformEndpointResult.EndpointArn, nil
	}
}

func (this *BasicSNSClient) PublishNotification(arn string, title string, message string, applicationType string) (err error) {

	callMeter, errorMeter := metrics.GetCallMeters("sns.publish_notification")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "Publish")
	values.Set("Message", formatMessage(applicationType, title, message))
	values.Set("MessageStructure", "json")
	values.Set("TargetArn", arn)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest("http://sns."+this.Region+".amazonaws.com/",
		values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
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

	callMeter, errorMeter := metrics.GetCallMeters("sns.make_request")
	callMeter.Mark(1)

	url_, err := url.Parse(host)
	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	aws.SignRequest(req, userId, userSecret, "sns", region)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

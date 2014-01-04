package sqs

import (
	"encoding/xml"
	"errors"
	"github.com/gopns/gopns/com/techtraits/gopns/aws"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Initilize(
	userId string,
	userSecret string,
	region string,
	queueName string) (error, *SqsQueue) {

	return CreateQueue(userId, userSecret, region, queueName)

}

func CreateQueue(userId string, userSecret string, region string, queueName string) (error, *SqsQueue) {
	values := url.Values{}
	values.Set("Action", "CreateQueue")
	values.Set("Version", "2012-11-05")
	values.Set("QueueName", queueName)
	values.Set("Attribute.1.Name", "ReceiveMessageWaitTimeSeconds")
	values.Set("Attribute.1.Value", "20")

	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest("http://sqs."+region+".amazonaws.com/",
		values, userId, userSecret, region)

	if err != nil {
		return err, nil
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		log.Printf("SQS Response %s", content)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message), nil
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		sqsQueue := new(SqsQueue)
		xml.Unmarshal(content, &sqsQueue)
		return nil, sqsQueue
	}

}

func makeRequest(host string, values url.Values, userId string,
	userSecret string, region string) (*http.Response, error) {

	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	aws.SignRequest(req, userId, userSecret, "sqs", region)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

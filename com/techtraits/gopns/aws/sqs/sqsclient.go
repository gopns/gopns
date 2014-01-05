package sqs

import (
	"encoding/xml"
	"errors"
	"github.com/gopns/gopns/com/techtraits/gopns/aws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

func SendMessage(userId string, userSecret string, region string, queueUrl string, message string) (error, *SqsSendMessageResponse) {
	values := url.Values{}
	values.Set("Action", "SendMessage")
	values.Set("Version", "2012-11-05")
	values.Set("MessageBody", message)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest(queueUrl, values, userId, userSecret, region)

	if err != nil {
		return err, nil
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message), nil
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		msgResp := new(SqsSendMessageResponse)
		xml.Unmarshal(content, &msgResp)
		return nil, msgResp
	}

}

func GetMessage(
	userId string,
	userSecret string,
	region string,
	queueUrl string,
	messageLimit int,
	waitTimeSeconds int) (error, []SqsMessage) {

	values := url.Values{}
	values.Set("Action", "ReceiveMessage")
	values.Set("Version", "2012-11-05")
	values.Set("AttributeName", "All")
	values.Set("MaxNumberOfMessages", strconv.Itoa(messageLimit))
	values.Set("WaitTimeSeconds", strconv.Itoa(waitTimeSeconds))
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest(queueUrl, values, userId, userSecret, region)

	if err != nil {
		return err, nil
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message), nil
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var msgResp SqsReceiveMessageResponse
		xml.Unmarshal(content, &msgResp)
		return nil, msgResp.SqsMessages
	}

}

func DeleteMessage(
	userId string,
	userSecret string,
	region string,
	queueUrl string,
	receiptHandle string) error {

	values := url.Values{}
	values.Set("Action", "DeleteMessage")
	values.Set("Version", "2012-11-05")
	values.Set("ReceiptHandle", receiptHandle)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := makeRequest(queueUrl, values, userId, userSecret, region)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		return nil
	}
}

func DeleteMessages(
	userId string,
	userSecret string,
	region string,
	queueUrl string,
	messages []SqsMessage) (err error, deletedMessageIds []string, messagesinError []ErrorMessage) {

	values := url.Values{}
	values.Set("Action", "DeleteMessageBatch")
	values.Set("Version", "2012-11-05")
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	for cnt, message := range messages {
		values.Set("DeleteMessageBatchRequestEntry."+strconv.Itoa(cnt+1)+".Id", message.MessageId)
		values.Set("DeleteMessageBatchRequestEntry."+strconv.Itoa(cnt+1)+".ReceiptHandle", message.ReceiptHandle)
	}

	response, err := makeRequest(queueUrl, values, userId, userSecret, region)

	if err != nil {
		return err, nil, nil
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message), nil, nil
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var msgResp SqsDeleteMessagesResponse
		xml.Unmarshal(content, &msgResp)
		return nil, msgResp.DeletedMessageIds, msgResp.Errors
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

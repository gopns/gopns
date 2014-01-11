package sqs

import (
	"encoding/xml"
	"errors"
	"github.com/gopns/gopns/aws"
	"github.com/gopns/gopns/metrics"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SQSClient interface {
	CreateQueue(queueName string) (*SqsQueue, error)
	SendMessage(queueUrl string, message string) (*SqsSendMessageResponse, error)
	GetMessage(queueUrl string, messageLimit int, waitTimeSeconds int) ([]SqsMessage, error)
	DeleteMessage(queueUrl string, receiptHandle string) error
	DeleteMessages(queueUrl string, messages []SqsMessage) (deletedMessageIds []string, messagesinError []ErrorMessage, err error)
}

type BasicSQSClient struct {
	UserId     string
	UserSecret string
	Region     string
}

func New(
	userId string,
	userSecret string,
	region string) SQSClient {
	return &BasicSQSClient{userId, userSecret, region}

}

func (this *BasicSQSClient) CreateQueue(queueName string) (*SqsQueue, error) {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.create_queue")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "CreateQueue")
	values.Set("Version", "2012-11-05")
	values.Set("QueueName", queueName)
	values.Set("Attribute.1.Name", "ReceiveMessageWaitTimeSeconds")
	values.Set("Attribute.1.Value", "20")

	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := this.makeRequest("http://sqs."+this.Region+".amazonaws.com/",
		values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		sqsQueue := new(SqsQueue)
		xml.Unmarshal(content, &sqsQueue)
		return sqsQueue, nil
	}

}

func (this *BasicSQSClient) SendMessage(queueUrl string, message string) (*SqsSendMessageResponse, error) {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.send_message")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "SendMessage")
	values.Set("Version", "2012-11-05")
	values.Set("MessageBody", message)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := this.makeRequest(queueUrl, values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		msgResp := new(SqsSendMessageResponse)
		xml.Unmarshal(content, &msgResp)
		return msgResp, nil
	}

}

func (this *BasicSQSClient) GetMessage(
	queueUrl string,
	messageLimit int,
	waitTimeSeconds int) ([]SqsMessage, error) {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.get_message")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "ReceiveMessage")
	values.Set("Version", "2012-11-05")
	values.Set("AttributeName", "All")
	values.Set("MaxNumberOfMessages", strconv.Itoa(messageLimit))
	values.Set("WaitTimeSeconds", strconv.Itoa(waitTimeSeconds))
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := this.makeRequest(queueUrl, values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var msgResp SqsReceiveMessageResponse
		xml.Unmarshal(content, &msgResp)
		return msgResp.SqsMessages, nil
	}

}

func (this *BasicSQSClient) DeleteMessage(
	queueUrl string,
	receiptHandle string) error {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.delete_message")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "DeleteMessage")
	values.Set("Version", "2012-11-05")
	values.Set("ReceiptHandle", receiptHandle)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	response, err := this.makeRequest(queueUrl, values, this.UserId, this.UserSecret, this.Region)

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
		return errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		return nil
	}
}

func (this *BasicSQSClient) DeleteMessages(
	queueUrl string,
	messages []SqsMessage) (deletedMessageIds []string, messagesinError []ErrorMessage, err error) {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.delete_messages")
	callMeter.Mark(1)

	values := url.Values{}
	values.Set("Action", "DeleteMessageBatch")
	values.Set("Version", "2012-11-05")
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	for cnt, message := range messages {
		values.Set("DeleteMessageBatchRequestEntry."+strconv.Itoa(cnt+1)+".Id", message.MessageId)
		values.Set("DeleteMessageBatchRequestEntry."+strconv.Itoa(cnt+1)+".ReceiptHandle", message.ReceiptHandle)
	}

	response, err := this.makeRequest(queueUrl, values, this.UserId, this.UserSecret, this.Region)

	if err != nil {
		errorMeter.Mark(1)
		return nil, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, nil, errors.New("Unable to create queue. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		var msgResp SqsDeleteMessagesResponse
		xml.Unmarshal(content, &msgResp)
		return msgResp.DeletedMessageIds, msgResp.Errors, nil
	}
}

func (this *BasicSQSClient) makeRequest(host string, values url.Values, userId string,
	userSecret string, region string) (*http.Response, error) {

	callMeter, errorMeter := metrics.GetCallMeters("sqs.make_request")
	callMeter.Mark(1)

	url_, err := url.Parse(host)

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	aws.SignRequest(req, userId, userSecret, "sqs", region)
	response, err := http.DefaultClient.Do(req)
	return response, err
}

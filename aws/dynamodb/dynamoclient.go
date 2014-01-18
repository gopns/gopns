package dynamodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gopns/gopns/aws"
	"github.com/gopns/gopns/metrics"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type DynamoClient interface {
	FindTable(dynamoTable string) (bool, error)
	CreateTable(createTableRequest CreateTableRequest) error
	UpdateItem(updateItemRequest UpdateItemRequest) error
	GetItem(getItemRequest GetItemRequest) (map[string]Attribute, error)
	ScanForItems(scanRequest ScanRequest) (*ScanResponse, error)
}

type requestor func(string, string, string) (int, io.ReadCloser, error)

type BasicDynamoClient struct {
	UserId     string
	UserSecret string
	Region     string
	r          requestor
}

func New(
	userId string,
	userSecret string,
	region string) (DynamoClient, error) {
	c := &BasicDynamoClient{UserId: userId, UserSecret: userSecret, Region: region}
	c.r = c.makeRequest
	return c, nil

}

func _new(
	userId string,
	userSecret string,
	region string,
	r requestor) (DynamoClient, error) {
	return &BasicDynamoClient{UserId: userId, UserSecret: userSecret, Region: region, r: r}, nil

}

func (this *BasicDynamoClient) FindTable(dynamoTable string) (bool, error) {

	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.findtable")
	callMeter.Mark(1)

	status, response, err := this.r("http://dynamodb."+this.Region+".amazonaws.com/",
		"{}", "ListTables")

	if err != nil {
		errorMeter.Mark(1)
		return false, err
	}

	defer response.Close()

	if status != 200 {
		content, _ := ioutil.ReadAll(response)
		var errorResponse aws.ErrorStruct
		err := json.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		if err != nil {
			return false, errors.New(fmt.Sprintf("Unable to serialize aws response. %s", content))
		} else {
			return false, errors.New("Unable to find table. " + errorResponse.Type + ": " + errorResponse.Message)
		}
	} else {

		content, _ := ioutil.ReadAll(response)
		var tableNames = make(map[string][]string)
		err := json.Unmarshal(content, &tableNames)
		if err != nil {
			return false, errors.New(fmt.Sprintf("Unable to serialize aws response.  %s", content))
		} else {
			for _, tableName := range tableNames["TableNames"] {
				if tableName == dynamoTable {
					log.Printf("Found Dynamo Table %s", tableName)
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (this *BasicDynamoClient) CreateTable(createTableRequest CreateTableRequest) error {

	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.createtable")
	callMeter.Mark(1)

	query, err := json.Marshal(createTableRequest)
	if err != nil {
		errorMeter.Mark(1)
		return err
	}

	status, response, err := this.r("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "CreateTable")

	if err != nil {
		errorMeter.Mark(1)
		return err
	}

	defer response.Close()

	if status != 200 {
		content, _ := ioutil.ReadAll(response)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return errors.New("Unable to create table. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		log.Printf("Created Dynamo Table %s", createTableRequest.TableName)
	}

	return nil
}

func (this *BasicDynamoClient) UpdateItem(updateItemRequest UpdateItemRequest) error {

	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.updateitem")
	callMeter.Mark(1)
	query, err := json.Marshal(updateItemRequest)
	if err != nil {
		errorMeter.Mark(1)
		return err
	}

	status, response, err := this.r("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "UpdateItem")

	if err != nil {
		errorMeter.Mark(1)
		return err
	}

	defer response.Close()

	if status != 200 {
		content, _ := ioutil.ReadAll(response)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return errors.New("Unable to update device. " + errorResponse.Type + ": " + errorResponse.Message)
	}

	return nil
}

func (this *BasicDynamoClient) GetItem(getItemRequest GetItemRequest) (map[string]Attribute, error) {

	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.getitem")
	callMeter.Mark(1)
	query, err := json.Marshal(getItemRequest)
	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	status, response, err := this.r("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "GetItem")

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	defer response.Close()

	if status != 200 {
		content, _ := ioutil.ReadAll(response)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, errors.New("Unable to get item. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		content, _ := ioutil.ReadAll(response)
		items := make(map[string]map[string]Attribute)
		json.Unmarshal(content, &items)
		if len(items) > 0 {
			return items["Item"], nil
		} else {
			return nil, errors.New("Unable to get item. Not Found")
		}
	}

}

func (this *BasicDynamoClient) ScanForItems(scanRequest ScanRequest) (*ScanResponse, error) {

	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.scanforitem")
	callMeter.Mark(1)
	query, err := json.Marshal(scanRequest)
	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	status, response, err := this.r("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "Scan")

	if err != nil {
		errorMeter.Mark(1)
		return nil, err
	}

	defer response.Close()

	if status != 200 {
		content, _ := ioutil.ReadAll(response)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		errorMeter.Mark(1)
		return nil, errors.New("Unable to scan for items. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		content, _ := ioutil.ReadAll(response)
		scanResponse := new(ScanResponse)
		json.Unmarshal(content, scanResponse)

		return scanResponse, nil
	}

}

func (this *BasicDynamoClient) makeRequest(host string, query string, action string) (int, io.ReadCloser, error) {
	callMeter, errorMeter := metrics.GetCallMeters("dynamodb.makerequest")
	callMeter.Mark(1)
	url_, err := url.Parse(host)
	if err != nil {
		errorMeter.Mark(1)
		return 0, nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(query))
	req.Header.Set("x-amz-target", "DynamoDB_20120810."+action)
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	aws.SignRequest(req, this.UserId, this.UserSecret, "dynamodb", this.Region)
	response, err := http.DefaultClient.Do(req)

	return response.StatusCode, response.Body, err
}

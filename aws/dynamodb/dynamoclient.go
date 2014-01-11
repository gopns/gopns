package dynamodb

import (
	"encoding/json"
	"errors"
	"github.com/gopns/gopns/aws"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type DynamoClientStruct struct {
	UserId     string
	UserSecret string
	Region     string
}

type DynamoClient interface {
	FindTable(dynamoTable string) (bool, error)
	CreateTable(createTableRequest CreateTableRequest) error
	UpdateItem(updateItemRequest UpdateItemRequest) error
	GetItem(getItemRequest GetItemRequest) (map[string]Attribute, error)
	ScanForItems(scanRequest ScanRequest) (*ScanResponse, error)
}

func Initilize(
	userId string,
	userSecret string,
	region string) (DynamoClient, error) {

	return &DynamoClientStruct{UserId: userId, UserSecret: userSecret, Region: region}, nil

}

func (this *DynamoClientStruct) FindTable(dynamoTable string) (bool, error) {
	response, err := this.makeRequest("http://dynamodb."+this.Region+".amazonaws.com/",
		"{}", "ListTables")

	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		return false, errors.New("Unable to register device. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {

		content, _ := ioutil.ReadAll(response.Body)
		var tableNames = make(map[string][]string)
		json.Unmarshal(content, &tableNames)
		for _, tableName := range tableNames["TableNames"] {
			if tableName == dynamoTable {
				log.Printf("Found Dynamo Table %s", tableName)
				return true, nil
			}
		}
	}

	return false, nil
}

func (this *DynamoClientStruct) CreateTable(createTableRequest CreateTableRequest) error {

	query, err := json.Marshal(createTableRequest)
	if err != nil {
		return err
	}

	response, err := this.makeRequest("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "CreateTable")

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		return errors.New("Unable to register device. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		log.Printf("Created Dynamo Table %s", createTableRequest.TableName)
	}

	return nil
}

func (this *DynamoClientStruct) UpdateItem(updateItemRequest UpdateItemRequest) error {

	query, err := json.Marshal(updateItemRequest)
	if err != nil {
		return err
	}

	response, err := this.makeRequest("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "UpdateItem")

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		return errors.New("Unable to register device. " + errorResponse.Type + ": " + errorResponse.Message)
	}

	return nil
}

func (this *DynamoClientStruct) GetItem(getItemRequest GetItemRequest) (map[string]Attribute, error) {

	query, err := json.Marshal(getItemRequest)
	if err != nil {
		return nil, err
	}

	response, err := this.makeRequest("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "GetItem")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		return nil, errors.New("Unable to register device. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		items := make(map[string]map[string]Attribute)
		json.Unmarshal(content, &items)
		return items["Item"], nil
	}

}

func (this *DynamoClientStruct) ScanForItems(scanRequest ScanRequest) (*ScanResponse, error) {

	query, err := json.Marshal(scanRequest)
	if err != nil {
		return nil, err
	}

	response, err := this.makeRequest("http://dynamodb."+this.Region+".amazonaws.com/",
		string(query[:]), "Scan")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorStruct
		json.Unmarshal(content, &errorResponse)
		return nil, errors.New("Unable to register device. " + errorResponse.Type + ": " + errorResponse.Message)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		scanResponse := new(ScanResponse)
		json.Unmarshal(content, scanResponse)

		return scanResponse, nil
	}

}

func (this *DynamoClientStruct) makeRequest(host string, query string, action string) (*http.Response, error) {

	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(query))
	req.Header.Set("x-amz-target", "DynamoDB_20120810."+action)
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	aws.SignRequest(req, this.UserId, this.UserSecret, "dynamodb", this.Region)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

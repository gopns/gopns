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

func Initilize(
	userId string,
	userSecret string,
	region string,
	dynamoTable string,
	initialReadCapacity int,
	initialWriteCapacity int) error {

	if found, err := findTable(userId, userSecret, region, dynamoTable); err != nil {
		return err
	} else if found {
		return nil
	} else {

		createTableRequest := CreateTableRequest{
			[]AttributeDefinition{AttributeDefinition{"alias", "S"}},
			dynamoTable,
			[]KeySchema{KeySchema{"alias", "HASH"}},
			ProvisionedThroughput{initialReadCapacity, initialWriteCapacity}}
		return createTable(userId, userSecret, region, createTableRequest)
	}

}

func findTable(userId string, userSecret string, region string, dynamoTable string) (bool, error) {
	response, err := makeRequest("http://dynamodb."+region+".amazonaws.com/",
		"{}", "ListTables", userId, userSecret, region)

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

func createTable(userId string, userSecret string, region string, createTableRequest CreateTableRequest) error {

	query, err := json.Marshal(createTableRequest)
	if err != nil {
		return err
	}

	response, err := makeRequest("http://dynamodb."+region+".amazonaws.com/",
		string(query[:]), "CreateTable", userId, userSecret, region)

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

func UpdateItem(updateItemRequest UpdateItemRequest, userId string,
	userSecert string, region string) error {

	query, err := json.Marshal(updateItemRequest)
	if err != nil {
		return err
	}

	response, err := makeRequest("http://dynamodb."+region+".amazonaws.com/",
		string(query[:]), "UpdateItem", userId, userSecert, region)

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

func GetItem(getItemRequest GetItemRequest, userId string,
	userSecert string, region string) (map[string]Attribute, error) {

	query, err := json.Marshal(getItemRequest)
	if err != nil {
		return nil, err
	}

	response, err := makeRequest("http://dynamodb."+region+".amazonaws.com/",
		string(query[:]), "GetItem", userId, userSecert, region)

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

func ScanForItems(scanRequest ScanRequest, userId string,
	userSecert string, region string) (*ScanResponse, error) {

	query, err := json.Marshal(scanRequest)
	if err != nil {
		return nil, err
	}

	response, err := makeRequest("http://dynamodb."+region+".amazonaws.com/",
		string(query[:]), "Scan", userId, userSecert, region)

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

func makeRequest(host string, query string, action string, userId string,
	userSecret string, region string) (*http.Response, error) {

	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(query))
	req.Header.Set("x-amz-target", "DynamoDB_20120810."+action)
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	aws.SignRequest(req, userId, userSecret, "dynamodb", region)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

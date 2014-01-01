package dynamodb

import (
	"encoding/json"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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

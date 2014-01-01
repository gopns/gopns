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

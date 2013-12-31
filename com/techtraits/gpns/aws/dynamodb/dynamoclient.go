package dynamodb

import (
	"encoding/json"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func UpdateItem(platformAppStr string, updateItemRequest UpdateItemRequest) error {

	platformApp := gpnsconfig.AWSConfigInstance().PlatformApps()[platformAppStr]
	updateItemRequest.TableName = platformApp.DynamoTable()

	query, err := json.Marshal(updateItemRequest)
	if err != nil {
		return err
	}

	response, err := MakeRequest("http://dynamodb."+platformApp.Region()+".amazonaws.com/",
		string(query[:]), platformAppStr)

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

func MakeRequest(host string, query string, platformAppName string) (*http.Response, error) {
	url_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url_.String(), strings.NewReader(query))
	req.Header.Set("x-amz-target", "DynamoDB_20120810.UpdateItem")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	aws.SignRequest(req, "dynamodb", platformAppName)
	response, err := http.DefaultClient.Do(req)

	return response, err
}

package sns

import (
	"encoding/xml"
	"errors"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PublisherStruct struct {
}

type Publisher interface {
	PublishNotification(platformAppName string, arn string, title string, message string) (err error)
}

func (this PublisherStruct) PublishNotification(platformAppName string, arn string, title string, message string) (err error) {
	values := url.Values{}
	values.Set("Action", "Publish")
	values.Set("Message", formatMessage(gpnsconfig.AWSConfigInstance().PlatformApps()[platformAppName].Type(), title, message))
	values.Set("MessageStructure", "json")
	values.Set("TargetArn", arn)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	url_, err := url.Parse("http://sns." + gpnsconfig.AWSConfigInstance().PlatformApps()[platformAppName].Region() + ".amazonaws.com/")
	if err != nil {
		return err
	}

	aws.SignRequest("POST", "/", values, url_.Host)

	response, err := http.PostForm(url_.String(), values)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		var errorResponse aws.ErrorResponse
		xml.Unmarshal(content, &errorResponse)
		return errors.New("Unable to send Push Notification. " + errorResponse.Error.Code + ": " + errorResponse.Error.Message)
	} else {
		return nil
	}

}

func formatMessage(appType string, title string, message string) string {
	formattedMessage := "{\"" + appType + "\":\"{\\\"data\\\":{\\\"id\\\":\\\"MESSAGE_ID\\\",\\\"title\\\":\\\"" + title + "\\\",\\\"alert\\\":\\\"" + message + "\\\"}}\"}"
	return formattedMessage
}

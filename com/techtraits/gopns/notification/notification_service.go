package notification

import (
	"code.google.com/p/gorest"
	"errors"
	"github.com/gopns/gopns/com/techtraits/gopns/aws/dynamodb"
	"github.com/gopns/gopns/com/techtraits/gopns/aws/sns"
	"github.com/gopns/gopns/com/techtraits/gopns/device"
	config "github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/rest/restutil"
	"strings"
)

type NotificationService struct {

	//Service level config
	gorest.RestService `root:"/rest/notification/" consumes:"application/json" produces:"application/json"`

	sendPushNotification gorest.EndPoint `method:"POST" path:"/{deviceAlias:string}" postdata:"NotificationMessage"`
	sendMassNotification gorest.EndPoint `method:"POST" path:"/?{localesParam:string}&{platformdParam:string}&{requiredTagsParam:string}&{skipTagsParam:string}" postdata:"NotificationMessage"`
}

func (serv NotificationService) SendPushNotification(message NotificationMessage, deviceAlias string) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	if !message.IsValid() {
		restutil.CheckError(errors.New("Invalid push notification message"), restError, 400)
	}

	key := make(map[string]dynamodb.Attribute)
	key["alias"] = dynamodb.Attribute{S: deviceAlias}
	getItemRequest := dynamodb.GetItemRequest{Key: key, TableName: config.AWSConfigInstance().DynamoTable()}

	item, err := dynamodb.GetItem(
		getItemRequest,
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())
	restutil.CheckError(err, restError, 500)

	if len(item) == 0 {
		restutil.CheckError(errors.New("Alias not "+deviceAlias+" not found"), restError, 404)
	} else {
		device_ := device.Device{item["alias"].S, item["locale"].S, item["arns"].SS, item["platform"].S, item["tags"].SS}
		for _, arn := range device_.Arns {

			sns.PublishNotification(
				arn,
				message.Title,
				message.Message,
				config.AWSConfigInstance().UserID(),
				config.AWSConfigInstance().UserSecret(),
				config.AWSConfigInstance().Region(),
				config.AWSConfigInstance().PlatformApps()[device_.Platform].Type())
		}
	}

}

func (serv NotificationService) SendMassNotification(
	message NotificationMessage,
	localesParam string,
	platformsParam string,
	requiredTagsParam string,
	skipTagsParam string) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

	err, _, _, _, _ := parseParameters(
		message,
		localesParam,
		platformsParam,
		requiredTagsParam,
		skipTagsParam)

	restutil.CheckError(err, restError, 400)
}

func parseParameters(message NotificationMessage, localesParam string, platformsParam string,
	requiredTagsParam string, skipTagsParam string) (error, []string, []string, []string, []string) {

	if !message.IsValid() {
		return errors.New("Invalid push notification message"), nil, nil, nil, nil
	}

	locales := strings.Split(localesParam, ",")
	platforms := strings.Split(platformsParam, ",")
	requiredTags := strings.Split(requiredTagsParam, ",")
	skipTags := strings.Split(skipTagsParam, ",")

	for _, locale := range locales {
		if err := device.ValidateLocale(locale); err != nil {
			return err, nil, nil, nil, nil
		}
	}

	for _, platform := range platforms {
		if err := device.ValidatePlatform(platform); err != nil {
			return err, nil, nil, nil, nil
		}
	}

	return nil, locales, platforms, requiredTags, skipTags
}
package notification

import (
	"code.google.com/p/gorest"
	"errors"
	"github.com/gopns/gopns/com/techtraits/gopns/aws/dynamodb"
	"github.com/gopns/gopns/com/techtraits/gopns/aws/sns"
	"github.com/gopns/gopns/com/techtraits/gopns/device"
	config "github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/rest/restutil"
)

type NotificationService struct {

	//Service level config
	gorest.RestService `root:"/rest/notification/" consumes:"application/json" produces:"application/json"`

	sendPushNotification gorest.EndPoint `method:"POST" path:"/{deviceAlias:string}" postdata:"Message"`
	sendMassNotification gorest.EndPoint `method:"POST" path:"/?{locale:string}&{platform:string}&{requiredTags:string}&{skipTags:string}" postdata:"Message"`
}

func (serv NotificationService) SendPushNotification(message Message, deviceAlias string) {

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
	message Message,
	locale string,
	platform string,
	requiredTags string,
	skipTags string) {

	restError := restutil.GetRestError(serv.ResponseBuilder())
	defer restutil.HandleErrors(restError)

}

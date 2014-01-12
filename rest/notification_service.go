package rest

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/device"
	"github.com/gopns/gopns/notification"
	"net/http"
)

type NotificationService struct {
	NotificationSender *notification.NotificationSender
	DeviceManager      device.DeviceManager
}

func (serv *NotificationService) Register(container *restful.Container, rootPath string) {
	ws := new(restful.WebService)
	ws.
		Path(rootPath + "/notification").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	ws.Route(ws.POST("/{deviceAlias}").To(serv.sendPushNotification).
		// docs
		Doc("send a push notification to the device by alias").
		Param(ws.PathParameter("deviceAlias", "the registered device alias").DataType("string")).
		Reads(notification.NotificationMessage{})) // from the request

	/*
		ws.Route(ws.POST("/").To(serv.sendMassNotification).
			// docs
			Doc("send mass push notifications to all users").
			Param(ws.QueryParameter("localesParam", "specify the locale").DataType("string")).
			Param(ws.QueryParameter("platformParam", "specify the platform").DataType("string")).
			Param(ws.QueryParameter("requiredTagsParam", "tags which must be set for a user, used for segmentation").DataType("string")).
			Param(ws.QueryParameter("skipTagsParam", "tags to skip for segmentation ").DataType("string")).
			Reads(notification.NotificationMessage{})) // from the request
	*/
	container.Add(ws)
}

func (serv *NotificationService) sendPushNotification(request *restful.Request, response *restful.Response) {

	alias := request.PathParameter("deviceAlias")
	err, device_ := serv.DeviceManager.GetDevice(alias)
	if err != nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Device not found.")
	}

	message := new(notification.NotificationMessage)
	err = request.ReadEntity(message)
	if err != nil {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}

	if !message.IsValid() {
		//ToDo use json error messages and appropriate error handling
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Invalid push notification message")
	}

	serv.NotificationSender.SendSyncNotification(*device_, *message, 5)

}

/*
func (serv NotificationService) sendMassNotification(
	message notification.NotificationMessage,
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

func parseParameters(message notification.NotificationMessage, localesParam string, platformsParam string,
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
*/

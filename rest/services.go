package rest

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/device"
	"github.com/gopns/gopns/gopnsconfig"
	"net/http"
)

func SetupRestServices() {

	gorest.RegisterService(new(device.DeviceService))
	gorest.RegisterService(new(NotificationService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+gopnsconfig.BaseConfigInstance().Port(), nil)
}

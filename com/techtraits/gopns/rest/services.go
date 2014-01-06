package rest

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/com/techtraits/gopns/device"
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"net/http"
)

func SetupRestServices() {

	gorest.RegisterService(new(device.DeviceService))
	gorest.RegisterService(new(NotificationService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+gopnsconfig.BaseConfigInstance().Port(), nil)
}

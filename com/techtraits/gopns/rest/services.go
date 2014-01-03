package rest

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/com/techtraits/gopns/device"
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"net/http"
)

func SetupRestServices() {

	deviceService := new(device.DeviceService)

	gorest.RegisterService(deviceService)
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+gopnsconfig.BaseConfigInstance().Port(), nil)
}

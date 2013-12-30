package rest

import (
	"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/device"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"net/http"
)

func SetupRestServices() {

	deviceService := new(device.DeviceService)

	gorest.RegisterService(deviceService)
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+gpnsconfig.BaseConfigInstance().Port(), nil)
}

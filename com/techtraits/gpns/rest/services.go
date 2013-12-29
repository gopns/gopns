package rest

import (
	"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/device"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"net/http"
)

func SetupRestServices(baseConfig gpnsconfig.BaseConfig) {
	gorest.RegisterService(new(device.DeviceService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+baseConfig.Port(), nil)
}

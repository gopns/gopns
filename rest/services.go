package rest

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/gopnsconfig"
	"net/http"
)

func SetupRestServices() {

	gorest.RegisterService(new(DeviceService))
	gorest.RegisterService(new(NotificationService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+gopnsconfig.BaseConfigInstance().Port(), nil)
}

package main

import (
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/rest"
)

func main() {

	appMode := gopnsconfig.ParseConfig()
	if appMode == gopnsconfig.SERVER_MODE {
		rest.SetupRestServices()
	} else if appMode == gopnsconfig.REGISTER_MODE {

	} else if appMode == gopnsconfig.SEND_MODE {

	}

}

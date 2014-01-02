package main

import (
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"github.com/usmanismail/gpns/com/techtraits/gpns/rest"
)

func main() {

	appMode := gpnsconfig.ParseConfig()
	if appMode == gpnsconfig.SERVER_MODE {
		rest.SetupRestServices()
	} else if appMode == gpnsconfig.REGISTER_MODE {

	} else if appMode == gpnsconfig.SEND_MODE {

	}

}

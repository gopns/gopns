package main

import (
	"github.com/gopns/gopns/gopnsapp"
	"github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/rest"
	"log"
)

func main() {

	appMode := gopnsconfig.ParseConfig()

	// start gopns -- MUST BE STARTED BEFORE ANYTHING ELSE
	gopnsapp_, err := gopnsapp.Initilize()
	if err == nil {
		gopnsapp_.Start()

		if appMode == gopnsconfig.SERVER_MODE {
			rest.SetupRestServices()
		} else if appMode == gopnsconfig.REGISTER_MODE {

		} else if appMode == gopnsconfig.SEND_MODE {

		}
	} else {
		log.Fatal(err.Error())
	}

}

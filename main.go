package main

import (
	"github.com/gopns/gopns/com/techtraits/gopns/aws/dynamodb"
	"github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/rest"
)

func main() {

	appMode := gopnsconfig.ParseConfig()
	dynamodb.Initilize(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().DynamoTable(),
		gopnsconfig.AWSConfigInstance().InitialReadCapacity(),
		gopnsconfig.AWSConfigInstance().InitialWriteCapacity())
	if appMode == gopnsconfig.SERVER_MODE {
		rest.SetupRestServices()
	} else if appMode == gopnsconfig.REGISTER_MODE {

	} else if appMode == gopnsconfig.SEND_MODE {

	}

}

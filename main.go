package main

import (
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"github.com/usmanismail/gpns/com/techtraits/gpns/rest"
)

func main() {

	gpnsconfig.ParseConfig()
	rest.SetupRestServices()

	/*


		publisher := sns.InitilizePublisher(awsConfig)
		arn, err := registrar.RegisterDevice("Test", "APA91bFaRKjCZfNcAhPTw6wSFGUxRi3108G_Swnz0fZ-Xr2pK9bwMGBjntXEJ72nrIyodMNx49cO3KESBpM3Jmd0zMpHToo1Cb_zR-_Lzqt5B-GRnzx3UuRHL6D6G9xaQwLLn05ugPjMm5Z8fLSTWocwT9ozCANcrqdM4tG-ljf7N3H7iSeymvo", "EN_US")
		if err != nil {
			log.Fatalf("Unable to register device %s", err.Error())
		}
		err = publisher.PublishNotification("Test", arn, "Title From GO", "Text From Go")
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	*/

}

package main

import (
	//"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws/sns"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"log"
	//"net/http"
)

func main() {

	baseConfig, awsConfig := gpnsconfig.ParseConfig()
	log.Printf("Running server on port %s", baseConfig.Port())
	log.Printf("Using AWS User ID %s", awsConfig.UserID())

	registrar := sns.Initilize(awsConfig)
	registrar.RegisterDevice("Test", "APA91bF1felnMnAgGtJm4NWcp2Zv4zpeKDko742sSdhBfK9uFtYREcoFQnLBuGockhSxMHMqTf2t5y_HwYe32PYVJNg0rwvGpdMbJwedgOZVdQ2lcQl6yB6CCp1xw2SosQcxU5JvGJLiO3aPuh53Qu_3Gzz-zpUgja2ZgLe31TAtHpY3Kgo3Fmc", "ENG_GB")

	/*
		gorest.RegisterService(new(RegistrationService))
		http.Handle("/", gorest.Handle())
		http.ListenAndServe(":8080", nil)*/

}

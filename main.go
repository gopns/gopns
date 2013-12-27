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
	log.Printf("Using AWS User Secret %s", awsConfig.UserSecret())

	registrar := sns.Initilize(awsConfig)
	registrar.RegisterDevice()

	/*
		gorest.RegisterService(new(RegistrationService))
		http.Handle("/", gorest.Handle())
		http.ListenAndServe(":8080", nil)*/

}

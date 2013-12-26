package main

import (
	"code.google.com/p/gorest"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"log"
	"net/http"
)

func main() {

	baseConfig, awsConfig := gpnsconfig.ParseConfig()
	log.Printf("Running server on port %s", baseConfig.Port())
	log.Printf("Using AWS User %s", awsConfig.UserID())

	gorest.RegisterService(new(RegistrationService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":8080", nil)

}

type RegistrationService struct {
	//Service level config
	gorest.RestService `root:"/device/" consumes:"application/json" produces:"application/json"`

	//End-Point level configs: Field names must be the same as the corresponding method names,
	// but not-exported (starts with lowercase)

	registerDevice gorest.EndPoint `method:"GET" path:"/" output:"string"`
}

func (serv RegistrationService) RegisterDevice() (arn string) {

	arn = "Test arn"
	return
}

package main

import (
	"github.com/gopns/gopns/gopnsapp"
	"log"
)

func main() {

	gopnsapp_, err := gopnsapp.New()
	if err == nil {
		gopnsapp_.Start()
	} else {
		log.Fatal(err.Error())
	}
}

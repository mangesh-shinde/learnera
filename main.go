package main

import (
	"log"
	"net/http"

	"github.com/mangesh-shinde/learnera/router"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":5000", r))
}

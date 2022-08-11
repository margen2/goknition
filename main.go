package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/margen2/goknition/config"
)

func main() {
	config.Load()
	fmt.Printf("Listening on localhost:%d", config.Port)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Port), nil))
}

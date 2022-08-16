package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/margen2/goknition/config"
	"github.com/margen2/goknition/controllers"
)

func main() {
	config.Load()
	fmt.Printf("Listening on localhost:%d\n", config.Port)
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/images", controllers.SearchImages)
	http.HandleFunc("/query", controllers.GetMatches)
	http.HandleFunc("/create-collection", controllers.CreateCollection)
	http.HandleFunc("/list-collections", controllers.GetCollections)
	http.HandleFunc("/delete-collection", controllers.DeleteCollection)
	http.HandleFunc("/pricing", controllers.Pricing)
	http.HandleFunc("/no-match", controllers.GetNoMatches)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Port), nil))
}

package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/margen2/goknition/config"
	"github.com/margen2/goknition/controllers"
)

//go:embed templates templates/partials templates/static
var static embed.FS

func main() {
	config.Load()
	fmt.Printf("Listening on localhost:%d\n", config.Port)

	controllers.LoadStatic(static)

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/create-collection", controllers.CreateCollection)
	http.HandleFunc("/list-collections", controllers.GetCollections)
	http.HandleFunc("/delete-collection", controllers.DeleteCollection)
	http.HandleFunc("/images", controllers.SearchImages)
	http.HandleFunc("/get-matches", controllers.GetMatches)
	http.HandleFunc("/no-match", controllers.GetNoMatches)
	http.HandleFunc("/save-matches", controllers.SaveMatches)
	http.HandleFunc("/pricing", controllers.Pricing)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/templates/", http.FileServer(http.FS(static)))

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Port), nil))
}

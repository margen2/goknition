package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/margen2/goknition/data"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.gohtml"))
}
func main() {
	fmt.Println("Listening on localhost:8080")
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		return
	}

	var images []data.Image
	images = data.LoadImages(r.FormValue("files"), images)
	IDs := data.LoadIDs(r.FormValue("ids"))

	for _, v := range images {
		fmt.Println(v.Path)
	}

	fmt.Println(strings.Repeat("-", 45))

	for _, v := range IDs {
		fmt.Println(v.Path)
	}
}

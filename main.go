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

	IDPath := r.FormValue("ids")
	filesPath := r.FormValue("pics")

	IDsData := data.LoadIDs(IDPath)
	filesData := data.LoadImages(filesPath)

	fmt.Println(IDsData)
	fmt.Println(strings.Repeat("-", 45))
	fmt.Println(filesData)
}

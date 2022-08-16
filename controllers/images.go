package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/margen2/goknition/api"
	"github.com/margen2/goknition/data"
	"github.com/margen2/goknition/db"
	"github.com/margen2/goknition/repositories"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SearchImages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "images.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	images, err := data.Load(r.FormValue("files"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)

	for _, image := range images {
		_, err := repositorie.CreateImage(image)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	collectionID := r.FormValue("collection")
	matches, nomatches, err := api.GetMatches(images, collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, match := range matches {
		for _, ID := range match.FaceIDs {
			err = repositorie.CreateMatch(ID, match.Image.FileName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	for _, nomatch := range nomatches {
		err := repositorie.CreateNoMatch(nomatch.FileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "query", http.StatusSeeOther)
}

func GetMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "query.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	query := r.FormValue("q")

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)
	images, err := repositorie.GetMatches(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "query.gohtml", images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetNoMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)
	images, err := repositorie.GetNoMatches()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "nomatch.gohtml", images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Pricing(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "pricing.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	faces, err := strconv.Atoi(r.FormValue("faces"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	images, err := strconv.Atoi(r.FormValue("images"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	facesn := fmt.Sprintf("Faces: %d * 0.001 = $%.2f ", faces, float64(faces)*0.001)
	imagesn := fmt.Sprintf("Images: %d * 0.001 = $%.2f ", images, float64(images)*0.001)
	cost := fmt.Sprintf("Total: %.2f + %.2f = $%.2f",
		float64(faces)*0.001, float64(images)*0.001, (float64(faces)*0.001)+(float64(images)*0.001))

	result := struct {
		Faces  string
		Images string
		Cost   string
	}{
		facesn,
		imagesn,
		cost,
	}

	err = tpl.ExecuteTemplate(w, "pricing.gohtml", result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

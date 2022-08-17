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

//Index serves all the routes that aren't already being served by another function.
//It serves the menu for Goknition.
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

//SearchImages serves the images route. It uploads all the files in the specified images
//folder to the Rekognition API and stores the result in the database.
func SearchImages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "images.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	images, err := data.Loadimages(r.FormValue("files"))
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
	http.Redirect(w, r, "get-matches", http.StatusSeeOther)
}

//Getmatches serves the get-matches route to return all of the images based on the given Face ID
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

//GetNoMatches serves the no-match route to return all of the images without a matching face.
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

//Pricing serves the pricing route to give an estimate on the pricing based on the amount of requests.
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

//SaveMatches serves the save-matches route to save all of the image matches on the given path.
func SaveMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "save-matches.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repositorie := repositories.NewImagesRepositorie(db)
	faces, err := repositorie.GetFaceIDs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("path")
	for _, face := range faces {
		images, err := repositorie.GetMatches(face.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = data.CopyImages(face.ID, path, images)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tpl.ExecuteTemplate(w, "save-matches.gohtml", path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

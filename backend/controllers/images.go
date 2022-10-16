package controllers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/margen2/goknition/backend/api"
	"github.com/margen2/goknition/backend/data"
	"github.com/margen2/goknition/backend/db"
	"github.com/margen2/goknition/backend/repositories"
)

var (
	tpl *template.Template
)

func LoadStatic(files embed.FS) {
	tpl = template.Must(template.ParseFS(files, "templates/*.html",
		"templates/partials/*.html", "templates/static/css/*.css", "templates/static/js/*.js"))
}

//Index serves all the routes that aren't already being served by another function.
//It serves the menu for Goknition.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//SearchImages serves the images route. It uploads all the files in the specified images
//folder to the Rekognition API and stores the result in the database.
func SearchImages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "images.html", nil)
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
	collection := r.FormValue("collection")
	collectionID, err := repositorie.GetCollectionID(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, image := range images {
		ID, err := repositorie.CreateImage(image, int(collectionID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images[i].ID = ID
	}

	matches, nomatches, err := api.GetMatches(images, collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, match := range matches {
		for _, ID := range match.FaceIDs {
			err = repositorie.CreateMatch(ID, match.Image.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	for _, nomatch := range nomatches {
		err := repositorie.CreateNoMatch(nomatch.ID)
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
		err := tpl.ExecuteTemplate(w, "get-matches.html", nil)
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

	err = tpl.ExecuteTemplate(w, "get-matches.html", images)
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

	err = tpl.ExecuteTemplate(w, "no-match.html", images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Pricing serves the pricing route to give an estimate on the pricing based on the amount of requests.
func Pricing(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "pricing.html", nil)
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

	facesn := fmt.Sprintf("Faces: %d * 0.001 = $%.3f ", faces, float64(faces)*0.001)
	imagesn := fmt.Sprintf("Images: %d * 0.001 = $%.3f ", images, float64(images)*0.001)
	cost := fmt.Sprintf("Total: %.3f + %.3f = $%.3f",
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

	err = tpl.ExecuteTemplate(w, "pricing.html", result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//SaveMatches serves the save-matches route to save all of the image matches on the given path.
func SaveMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "save-matches.html", nil)
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
	collection := r.FormValue("collection")
	collectionID, err := repositorie.GetCollectionID(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	faces, err := repositorie.GetFaceIDs(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("path")
	for _, face := range faces {
		fmt.Println("printin faces", face.FaceID)
		images, err := repositorie.GetMatches(face.FaceID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = data.CopyImages(face.FaceID, path, images)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tpl.ExecuteTemplate(w, "save-matches.html", path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

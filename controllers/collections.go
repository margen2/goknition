package controllers

import (
	"net/http"

	"github.com/margen2/goknition/api"
	"github.com/margen2/goknition/data"
	"github.com/margen2/goknition/db"
	"github.com/margen2/goknition/repositories"
)

//CreateCollection serves the create-collection route. It create a new Rekognition Collection
// and index all of the faces in the IDs folder to it.
func CreateCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "collections.gohtml", nil)
		if err != nil {
			return
		}
		return
	}

	faces, err := data.LoadFaces(r.FormValue("ids"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)

	for _, face := range faces {
		imageID, err := repositorie.CreateImage(face.Image)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = repositorie.CreateFace(face, imageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	collectionID := r.FormValue("collection")
	err = api.PrepareCollection(collectionID, faces)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "collections.gohtml", collectionID)
	if err != nil {
		return
	}
}

//GetCollections serves the get-collections route to return all of the active collections on the
// Rekognition APi
func GetCollections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	collections, err := api.ListCollections()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "list-collections.gohtml", collections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//DeleteCollection serves the delete-collection route to delete the
//requested collection from the Rekognition APi.
func DeleteCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "delete-collection.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	collectionID := r.FormValue("collection")
	err := api.DeleteCollection(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "delete-collection.gohtml", collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

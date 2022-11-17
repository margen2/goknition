package controllers

import (
	"github.com/margen2/goknition/backend/api"
	"github.com/margen2/goknition/backend/data"
	"github.com/margen2/goknition/backend/db"
	"github.com/margen2/goknition/backend/repositories"
)

//CreateCollection creates a new Rekognition Face Collection.
func CreateCollection(collection string) error {
	db, err := db.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)
	err = repositorie.CreateCollection(collection)
	if err != nil {
		return err
	}

	// err = api.CreateCollection(collection)
	// if err != nil {
	// 	return err
	// }

	return nil
}

//DeleteCollection deletes the given collection from the Rekognition API.
func DeleteCollection(collectionID string) error {
	err := api.DeleteCollection(collectionID)
	if err != nil {
		return err
	}

	return nil
}

//IndexFaces indexes the faces found in the path to the given collectionID.
func IndexFaces(collection, path string) error {
	faces, err := data.LoadFaces(path)
	if err != nil {
		return err
	}

	db, err := db.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)

	collectionID, err := repositorie.GetCollectionID(collection)
	if err != nil {
		return err
	}

	for _, face := range faces {
		_, err := repositorie.CreateFace(face.FaceID, collectionID)
		if err != nil {
			return err
		}
	}

	err = api.IndexFaces(collection, faces)
	if err != nil {
		return err
	}

	return nil
}

//GetFaces returns the saved faces for the given collection.
func GetFaces(collection string) ([]string, error) {

	db, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)

	collectionID, err := repositorie.GetCollectionID(collection)
	if err != nil {
		return nil, err
	}

	faces, err := repositorie.GetFaces(collectionID)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, face := range faces {
		result = append(result, face.FaceID)
	}

	return result, nil
}

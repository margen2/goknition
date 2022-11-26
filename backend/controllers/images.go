package controllers

import (
	"github.com/margen2/goknition/backend/api"
	"github.com/margen2/goknition/backend/data"
	"github.com/margen2/goknition/backend/db"
	"github.com/margen2/goknition/backend/models"
	"github.com/margen2/goknition/backend/repositories"
)

//SearchImages uploads all the files in the given images
//folder to the Rekognition API and saves the result to the database.
func SearchImages(collection, path string) error {
	images, err := data.Loadimages(path)
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

	for i, image := range images {
		ID, err := repositorie.CreateImage(image, collectionID)
		if err != nil {
			return err
		}
		images[i].ID = ID
	}

	matches, nomatches, err := api.GetMatches(images, collection)
	if err != nil {
		return err
	}

	faces, err := repositorie.GetFaceIDs(collectionID)
	if err != nil {
		return err
	}
	faceIDs := make(map[string]uint64, len(faces))
	for _, face := range faces {
		faceIDs[face.FaceID] = face.ID
	}

	for _, match := range matches {
		for _, ID := range match.FaceIDs {
			err = repositorie.CreateMatch(faceIDs[ID], match.Image.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, nomatch := range nomatches {
		err := repositorie.CreateNoMatch(nomatch.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

//Getmatches returns all of the images based on the given Face ID
func GetMatches(FaceID string) ([]models.Image, error) {
	db, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)

	images, err := repositorie.GetMatches(FaceID)
	if err != nil {
		return nil, err
	}
	return images, nil
}

//GetNoMatches returns all the images without a matching face.
func GetNoMatches() ([]models.Image, error) {
	db, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repositorie := repositories.NewImagesRepositorie(db)
	images, err := repositorie.GetNoMatches()
	if err != nil {
		return nil, err
	}
	return images, nil
}

//SaveMatches saves all of the image matches on the given path.
func SaveMatches(collection, path string) error {
	db, err := db.ConnectDB()
	if err != nil {
		return err
	}

	repositorie := repositories.NewImagesRepositorie(db)
	collectionID, err := repositorie.GetCollectionID(collection)
	if err != nil {
		return err
	}

	faces, err := repositorie.GetFaceIDs(collectionID)
	if err != nil {
		return err
	}

	for _, face := range faces {
		images, err := repositorie.GetMatches(face.FaceID)

		if err != nil {
			return err
		}

		err = data.CopyImages(face.FaceID, path, images)
		if err != nil {
			return err
		}
	}
	return nil
}

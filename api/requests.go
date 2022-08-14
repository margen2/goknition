package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

// PrepareCollection creates a new rekognition collection and adds all of the given faces
// to it.
func PrepareCollection(collectionID string, faces []models.Face) error {
	svc := newClient()

	err := createCollection(svc, collectionID)
	if err != nil {
		return fmt.Errorf("createcollection: %w", err)
	}

	err = indexFaces(svc, collectionID, faces)
	if err != nil {
		return fmt.Errorf("indexfaces: %w", err)
	}

	return nil
}

// ListCollections returns a list that contains all of the collection IDs created by
// the connected AWS account.
func ListCollections() ([]string, error) {
	svc := newClient()
	input := &rekognition.ListCollectionsInput{}

	result, err := svc.ListCollections(input)
	if err != nil {
		return nil, fmt.Errorf("svc.listcollection: %w", err)
	}

	var collectionsIDs []string
	for _, ID := range result.CollectionIds {
		collectionsIDs = append(collectionsIDs, *ID)
	}

	return collectionsIDs, nil
}

// Getmatches receives a list of images and returns all of the matches found using the specified collections.
func GetMatches(images []models.Image, collectionID string) ([]models.Match, []models.Image, error) {
	svc := newClient()

	var matches []models.Match
	var nomatches []models.Image
	l := len(images) - 1

	for i, image := range images {
		fmt.Printf("image N°%d out of %d\n", i, l)
		result, err := searchFaces(svc, image, collectionID)
		if err != nil {
			return nil, nil, fmt.Errorf("getmatches: %w", err)
		}

		match := models.Match{image, nil}
		if len(result.FaceMatches) > 0 {
			for _, fm := range result.FaceMatches {
				match.FaceIDs = append(match.FaceIDs, *fm.Face.ExternalImageId)
			}
		} else {
			nomatches = append(nomatches, image)
		}
		matches = append(matches, match)
	}

	return matches, nomatches, nil
}

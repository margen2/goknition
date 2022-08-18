package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

// Getmatches receives a list of images and returns all of the matches found using the specified collections.
func GetMatches(images []models.Image, collectionID string) ([]models.Match, []models.Image, error) {
	errorC := make(chan error)
	matchC := make(chan models.Match)
	noMatchC := make(chan models.Image)

	svc := newClient()
	for _, image := range images {
		go searchFaces(svc, image, collectionID, matchC, errorC, noMatchC)
	}
	fmt.Printf("started all of the  %d images\n", len(images)-1)

	l := len(images)
	var matches []models.Match
	var nomatches []models.Image
	for i := 0; i < l; i++ {
		select {
		case match := <-matchC:
			matches = append(matches, match)
		case noMatch := <-noMatchC:
			nomatches = append(nomatches, noMatch)
		case err := <-errorC:
			return nil, nil, err
		}
	}

	return matches, nomatches, nil
}

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

// DeleteCollection deletes the corresponding collection from AWS
func DeleteCollection(collectionID string) error {
	svc := newClient()

	input := &rekognition.DeleteCollectionInput{
		CollectionId: aws.String(collectionID),
	}

	_, err := svc.DeleteCollection(input)
	if err != nil {
		return fmt.Errorf("svc.deletecollection: %w", err)
	}

	return nil
}

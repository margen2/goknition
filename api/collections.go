package api

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

func createCollection(svc *rekognition.Rekognition, collectionId string) error {

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionId),
	}

	_, err := svc.CreateCollection(input)
	if err != nil {
		return fmt.Errorf("svc.createcollection: %w", err)
	}
	return nil
}

func indexFaces(svc *rekognition.Rekognition, collectionId string, faces []models.Face) error {

	for _, face := range faces {
		imageAWS, err := newImageAWS(face.Image.Path + `\` + face.Image.Filename)
		if err != nil {
			return fmt.Errorf("newimageaws: %w", err)
		}

		input := &rekognition.IndexFacesInput{
			CollectionId:    aws.String(collectionId),
			Image:           imageAWS,
			ExternalImageId: aws.String(strconv.Itoa(int(face.ID))),
			MaxFaces:        aws.Int64(1),
		}

		_, err = svc.IndexFaces(input)
		if err != nil {
			return fmt.Errorf("svc.indexfaces: %w", err)
		}
	}
	return nil
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

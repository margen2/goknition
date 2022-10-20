package api

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/backend/models"
)

var collectionIDs []string

// CreateCollection
func CreateCollection(collectionId string) error {
	svc := newClient()

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionId),
	}

	_, err := svc.CreateCollection(input)
	if err != nil {
		return fmt.Errorf("svc.createcollection: %w", err)
	}
	return nil
}

// IndexFaces
func IndexFaces(collectionId string, faces []models.Face) error {
	svc := newClient()

	for _, face := range faces {
		imageAWS, err := newImageAWS(filepath.Join(face.Image.Path, face.Image.Filename))
		if err != nil {
			return fmt.Errorf("newimageaws: %w", err)
		}

		input := &rekognition.IndexFacesInput{
			CollectionId:    aws.String(collectionId),
			Image:           imageAWS,
			ExternalImageId: aws.String(face.ID),
			MaxFaces:        aws.Int64(1),
		}

		_, err = svc.IndexFaces(input)
		if err != nil {
			return fmt.Errorf("svc.indexfaces: %w", err)
		}
	}
	return nil
}

// RefreshCollections returns a list that contains all collection IDs created by
// the connected AWS account.
func RefreshCollections() error {
	svc := newClient()
	input := &rekognition.ListCollectionsInput{}

	result, err := svc.ListCollections(input)
	if err != nil {
		return fmt.Errorf("svc.listcollection: %w", err)
	}

	collectionIDs = nil
	for _, ID := range result.CollectionIds {
		collectionIDs = append(collectionIDs, *ID)
	}

	return nil
}

// ListCollections lists
func ListCollections() []string {
	return collectionIDs
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

package api

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/backend/models"
)

// CreateCollection creates a new face collection on AWS
func CreateCollection(collectionId string) error {
	svc := rekognition.New(mySession)

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionId),
	}

	_, err := svc.CreateCollection(input)
	if err != nil {
		return fmt.Errorf("svc.createcollection: %w", err)
	}
	return nil
}

// IndexFaces adds the faces array to the given collectionID
func IndexFaces(collectionId string, faces []models.Face) error {
	svc := rekognition.New(mySession)

	for _, face := range faces {
		for _, image := range face.Images {
			imageAWS, err := newImageAWS(filepath.Join(image.Path, image.Filename))
			if err != nil {
				return fmt.Errorf("newimageaws: %w", err)
			}

			input := &rekognition.IndexFacesInput{
				CollectionId:    aws.String(collectionId),
				Image:           imageAWS,
				ExternalImageId: aws.String(face.FaceID),
				MaxFaces:        aws.Int64(1),
			}
			_, err = svc.IndexFaces(input)
			if err != nil {
				return fmt.Errorf("svc.indexfaces: %w", err)
			}
		}
	}
	return nil
}

// collectionIDs represents the active collections on AWS
var collectionIDs []string

// RefreshCollections updates the collectionsIDs variable
func RefreshCollections() error {
	svc := rekognition.New(mySession)

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

// ListCollections lists the active collections
func ListCollections() []string {
	return collectionIDs
}

// DeleteCollection deletes the corresponding collection from AWS
func DeleteCollection(collectionID string) error {
	svc := rekognition.New(mySession)

	input := &rekognition.DeleteCollectionInput{
		CollectionId: aws.String(collectionID),
	}

	_, err := svc.DeleteCollection(input)
	if err != nil {
		return fmt.Errorf("svc.deletecollection: %w", err)
	}

	return nil
}

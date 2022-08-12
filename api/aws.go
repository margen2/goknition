package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

func indexFaces(svc *rekognition.Rekognition, collectionId string, faces []models.Face) error {

	for _, face := range faces {
		imageAWS, err := newImageAWS(face.Image.Path + `\` + face.Image.FileName + ".JPG")
		if err != nil {
			return fmt.Errorf("indexfaces: %w", err)
		}

		input := &rekognition.IndexFacesInput{
			CollectionId:    aws.String(collectionId),
			Image:           imageAWS,
			ExternalImageId: aws.String(face.ID),
			MaxFaces:        aws.Int64(1),
		}

		_, err = svc.IndexFaces(input)
		if err != nil {
			return fmt.Errorf("indexcollection/svc.indexfaces: %w", err)
		}
	}
	return nil
}

func searchFaces(svc *rekognition.Rekognition, image models.Image, collectionId string) (*rekognition.SearchFacesByImageOutput, error) {

	imageAWS, err := newImageAWS(image.Path + `\` + image.FileName + ".JPG")
	if err != nil {
		return nil, fmt.Errorf("searchfaces: %w", err)
	}

	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(collectionId),
		FaceMatchThreshold: aws.Float64(95.000000),
		Image:              imageAWS,
		MaxFaces:           aws.Int64(5),
	}

	result, err := svc.SearchFacesByImage(input)
	if err != nil {
		return nil, fmt.Errorf("searchfaces/svc.searchfacesbyimage: %w", err)
	}
	return result, nil
}

func createCollection(svc *rekognition.Rekognition, collectionId string) error {

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionId),
	}

	_, err := svc.CreateCollection(input)
	if err != nil {
		return fmt.Errorf("createcollection: %w", err)
	}
	return nil
}

package api

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

func newImageAWS(path string) (*rekognition.Image, error) {
	b, err := getImageBytes(path)
	if err != nil {
		return nil, fmt.Errorf("newimageaws: %w", err)
	}
	return &rekognition.Image{
		Bytes: b,
	}, nil
}

func newClient() *rekognition.Rekognition {
	mySession := session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(2),
		Region:     aws.String("us-east-1"),
	}))

	return rekognition.New(mySession)
}

func getImageBytes(filePath string) ([]byte, error) {
	fl, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("getimagebytes/os.open: %w", err)
	}
	defer fl.Close()

	fileInfo, err := fl.Stat()
	if err != nil {
		return nil, fmt.Errorf("getimagebytes/fl.stat: %w", err)
	}

	b := make([]byte, fileInfo.Size())
	n, err := fl.Read(b)
	if err != nil || n == 0 {
		return nil, fmt.Errorf("getimagebytes/fl.read: %w", err)
	}

	return b, nil
}

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

func searchFaces(svc *rekognition.Rekognition, image models.Image, collectionId string,
	matchC chan models.Match, errC chan error, noMatchC chan models.Image) {

	imageAWS, err := newImageAWS(image.Path + `\` + image.FileName + ".JPG")
	if err != nil {
		errC <- fmt.Errorf("newimageaws: %w", err)
	}

	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(collectionId),
		FaceMatchThreshold: aws.Float64(95.000000),
		Image:              imageAWS,
		MaxFaces:           aws.Int64(5),
	}

	result, err := svc.SearchFacesByImage(input)
	if err != nil {
		errC <- fmt.Errorf("svc.searchfacesbyimage: %w", err)
	}

	match := models.Match{image, nil}
	if len(result.FaceMatches) > 0 {
		for _, fm := range result.FaceMatches {
			match.FaceIDs = append(match.FaceIDs, *fm.Face.ExternalImageId)
		}
		matchC <- match
	} else {
		noMatchC <- image
	}
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

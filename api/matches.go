package api

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

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
	var size int64 = fileInfo.Size()
	b := make([]byte, size)

	n, err := fl.Read(b)
	if err != nil || n == 0 {
		return nil, fmt.Errorf("getimagebytes/fl.read: %w", err)
	}

	return b, nil
}

func newImageAWS(path string) (*rekognition.Image, error) {
	b, err := getImageBytes(path)
	if err != nil {
		return nil, fmt.Errorf("newimageaws: %w", err)
	}
	return &rekognition.Image{
		Bytes: b,
	}, nil
}

func GetMatches(faces []models.Face, images []models.Image, collectionID string) ([]models.Match, []models.Image, error) {
	svc := newClient()

	err := createCollection(svc, collectionID)
	if err != nil {
		return nil, nil, fmt.Errorf("getmatches: %w", err)
	}

	err = indexFaces(svc, collectionID, faces)
	if err != nil {
		return nil, nil, fmt.Errorf("getmatches: %w", err)
	}

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

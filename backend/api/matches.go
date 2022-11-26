package api

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/backend/models"
)

type result struct {
	match   models.Match
	noMatch models.Image
}

type imageBytes struct {
	image models.Image
	bytes []byte
}

// Getmatches receives a list of images and returns all of the matches found using the specified collection.
func GetMatches(images []models.Image, collectionID string) ([]models.Match, []models.Image, error) {

	compressorC := make(chan imageBytes, 10)
	requestC := make(chan imageBytes)
	resultC := make(chan result)
	errorC := make(chan error)

	//50 is the default number of requests per second for the Rekognition API.
	//If you ask for a limit increase, you need to increase the number of workers accordingly.
	const workers = 50
	svc := rekognition.New(mySession)
	for i := 0; i < workers; i++ {
		go searchFaces(svc, requestC, resultC, errorC, collectionID)
	}

	for i := 0; i < runtime.NumCPU()/2; i++ {
		go imageCompressor(compressorC, requestC, errorC)
	}

	go sendImageBytes(images, compressorC, requestC, errorC)

	l := len(images)
	var matches []models.Match
	var nomatches []models.Image
	for i := 0; i < l; i++ {
		fmt.Printf("image N° %d out of %d\n", i+1, l)
		select {
		case res := <-resultC:
			if res.match.FaceIDs == nil {
				nomatches = append(nomatches, res.noMatch)
			} else {
				matches = append(matches, res.match)
			}
		case err := <-errorC:
			return nil, nil, err
		}
	}

	return matches, nomatches, nil
}

func sendImageBytes(images []models.Image, compressorC chan imageBytes, requestC chan imageBytes, errC chan error) {
	for _, image := range images {
		b, err := getImageBytes(filepath.Join(image.Path, image.Filename))
		if err != nil {
			errC <- err
			return
		}

		// 5 mb is the limit accepted by the rekognition API
		if len(b) > 5242880 {
			compressorC <- imageBytes{image, b}
			continue
		}

		requestC <- imageBytes{image, b}
	}
	close(compressorC)
}

func imageCompressor(compressorC chan imageBytes, requestC chan imageBytes, errC chan error) {
	for image := range compressorC {
		quality := 75
		for len(image.bytes) > 5242880 {
			var err error
			image.bytes, err = compressImage(image.bytes, quality)
			if err != nil {
				errC <- err
				return
			}
			quality -= 10
		}
		requestC <- image
	}
}

func searchFaces(svc *rekognition.Rekognition, requestC chan imageBytes, resultC chan result, errC chan error, collectionID string) {
	for r := range requestC {

		input := &rekognition.SearchFacesByImageInput{
			CollectionId:       aws.String(collectionID),
			FaceMatchThreshold: aws.Float64(75.00),
			Image:              &rekognition.Image{Bytes: r.bytes},
			MaxFaces:           aws.Int64(6),
		}

		res, err := svc.SearchFacesByImage(input)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				//do nothing
			default:
				errC <- fmt.Errorf("svc.searchfacesbyimage: %w - at %s", err, filepath.Join(r.image.Path, r.image.Filename))
				return
			}
		}

		match := models.Match{Image: r.image, FaceIDs: nil}
		if len(res.FaceMatches) > 0 {
			for _, fm := range res.FaceMatches {
				match.FaceIDs = append(match.FaceIDs, *fm.Face.ExternalImageId)
			}
			resultC <- result{match, models.Image{}}
		} else {
			resultC <- result{match, r.image}
		}
	}
}

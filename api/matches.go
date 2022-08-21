package api

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/models"
)

type request struct {
	image models.Image
	input *rekognition.SearchFacesByImageInput
}
type result struct {
	match   models.Match
	noMatch models.Image
}

// Getmatches receives a list of images and returns all of the matches found using the specified collection.
func GetMatches(images []models.Image, collectionID string) ([]models.Match, []models.Image, error) {
	errorC := make(chan error)
	requestC := make(chan request)
	resultC := make(chan result)

	l := len(images)
	for i := 0; i < l; i++ {
		svc := newClient()
		go searchFaces(svc, requestC, resultC, errorC)
	}

	for _, image := range images {
		imageAWS, err := newImageAWS(image.Path + `\` + image.Filename)
		if err != nil {
			return nil, nil, fmt.Errorf("newimageaws: %w", err)
		}

		input := &rekognition.SearchFacesByImageInput{
			CollectionId:       aws.String(collectionID),
			FaceMatchThreshold: aws.Float64(95.000000),
			Image:              imageAWS,
			MaxFaces:           aws.Int64(5),
		}

		requestC <- request{image, input}
	}

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
	close(requestC)

	return matches, nomatches, nil
}

func searchFaces(svc *rekognition.Rekognition, requestsC chan request, resultC chan result, errC chan error) {
	r := <-requestsC
	res, err := svc.SearchFacesByImage(r.input)

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case rekognition.ErrCodeInvalidParameterException:
			//do nothing
		default:
			errC <- fmt.Errorf("svc.searchfacesbyimage: %w", err)
			return
		}
	}

	match := models.Match{r.image, nil}
	if len(res.FaceMatches) > 0 {
		for _, fm := range res.FaceMatches {
			FaceID, err := strconv.ParseUint(*fm.Face.ExternalImageId, 10, 64)
			if err != nil {
				errC <- fmt.Errorf("strconv.atoi: %w", err)
				return
			}
			match.FaceIDs = append(match.FaceIDs, FaceID)
		}
		resultC <- result{match, models.Image{}}
	} else {
		resultC <- result{match, r.image}
	}
}

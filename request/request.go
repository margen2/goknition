package request

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/margen2/goknition/data"
)

func IndexCollection(svc *rekognition.Rekognition, collectionId string, faces []data.Face) {

	for _, face := range faces {
		imageAWS := newImageAWS(face.Image.Bytes)

		input := &rekognition.IndexFacesInput{
			CollectionId:    aws.String(collectionId),
			Image:           imageAWS,
			ExternalImageId: aws.String(face.ID),
			MaxFaces:        aws.Int64(1),
		}

		_, err := svc.IndexFaces(input)
		if err != nil {
			log.Println("indexfaces:")
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case rekognition.ErrCodeInvalidS3ObjectException:
					fmt.Println(rekognition.ErrCodeInvalidS3ObjectException, aerr.Error())
				case rekognition.ErrCodeInvalidParameterException:
					fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
				case rekognition.ErrCodeImageTooLargeException:
					fmt.Println(rekognition.ErrCodeImageTooLargeException, aerr.Error())
				case rekognition.ErrCodeAccessDeniedException:
					fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
				case rekognition.ErrCodeInternalServerError:
					fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
				case rekognition.ErrCodeThrottlingException:
					fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
				case rekognition.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case rekognition.ErrCodeResourceNotFoundException:
					fmt.Println(rekognition.ErrCodeResourceNotFoundException, aerr.Error())
				case rekognition.ErrCodeInvalidImageFormatException:
					fmt.Println(rekognition.ErrCodeInvalidImageFormatException, aerr.Error())
				case rekognition.ErrCodeServiceQuotaExceededException:
					fmt.Println(rekognition.ErrCodeServiceQuotaExceededException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			os.Exit(1)
		}
	}
}

func searchFaces(svc *rekognition.Rekognition, imageAWS *rekognition.Image, collectionId string) *rekognition.SearchFacesByImageOutput {
	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(collectionId),
		FaceMatchThreshold: aws.Float64(95.000000),
		Image:              imageAWS,
		MaxFaces:           aws.Int64(5),
	}

	result, err := svc.SearchFacesByImage(input)
	if err != nil {
		log.Println("searchfacesbyimage:")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidS3ObjectException:
				fmt.Println(rekognition.ErrCodeInvalidS3ObjectException, aerr.Error())
			case rekognition.ErrCodeInvalidParameterException:
				fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				fmt.Println(rekognition.ErrCodeImageTooLargeException, aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case rekognition.ErrCodeResourceNotFoundException:
				fmt.Println(rekognition.ErrCodeResourceNotFoundException, aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				fmt.Println(rekognition.ErrCodeInvalidImageFormatException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	return result
}

func CreateCollection(svc *rekognition.Rekognition, collectionId string) {

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionId),
	}

	_, err := svc.CreateCollection(input)
	if err != nil {
		log.Println("createcollection:")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case rekognition.ErrCodeResourceAlreadyExistsException:
				fmt.Println(rekognition.ErrCodeResourceAlreadyExistsException, aerr.Error())
			case rekognition.ErrCodeServiceQuotaExceededException:
				fmt.Println(rekognition.ErrCodeServiceQuotaExceededException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
}

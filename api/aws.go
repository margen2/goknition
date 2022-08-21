package api

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

func newImageAWS(path string) (*rekognition.Image, error) {
	b, err := getImageBytes(path)
	if err != nil {
		return nil, fmt.Errorf("getimagebytes: %w", err)
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
		return nil, fmt.Errorf("os.open: %w", err)
	}
	defer fl.Close()

	fileInfo, err := fl.Stat()
	if err != nil {
		return nil, fmt.Errorf("fl.stat: %w", err)
	}

	b := make([]byte, fileInfo.Size())
	n, err := fl.Read(b)
	if err != nil || n == 0 {
		return nil, fmt.Errorf("fl.read: %w", err)
	}

	return b, nil
}

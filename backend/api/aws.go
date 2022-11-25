package api

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/edwvee/exiffix"
)

var mySession *session.Session

func InitializeSession() {
	mySession = session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(2),
		Region:     aws.String("us-east-1"),
	}))
}

func newImageAWS(path string) (*rekognition.Image, error) {
	b, err := getImageBytes(path)
	if err != nil {
		return nil, fmt.Errorf("getimagebytes: %w", err)
	}
	return &rekognition.Image{
		Bytes: b,
	}, nil
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

// compressImage takes a slice of byte that represents a JPEG file and
// compresses it to match the given quality.
// The image/jpeg package is not being used to decode the image because it deletes the Exif
// data in the process, which may result in changes to the image orientation.
func compressImage(data []byte, quality int) ([]byte, error) {
	img, _, err := exiffix.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("exiffix.decode: %w", err)
	}

	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("jpeg.encode: %w", err)
	}

	return buf.Bytes(), nil
}

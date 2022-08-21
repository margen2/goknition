package models

// Face represents a face value within the given image.
type Face struct {
	ID     uint64
	FaceID string
	Image  Image
}

// Image represents a image file.
type Image struct {
	ID       uint64
	Filename string
	Path     string
}

// Match represents the matches between one image and all the faces that it contains.
type Match struct {
	Image   Image
	FaceIDs []uint64
}

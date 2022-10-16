package repositories

import (
	"database/sql"
	"fmt"

	"github.com/margen2/goknition/backend/models"
)

// Image represent a pointer to type sql.DB.
type Images struct {
	db *sql.DB
}

// NewImagesRepositorie returns a pointer to type Images that allows database methods.
func NewImagesRepositorie(db *sql.DB) *Images {
	return &Images{db}
}

func (repositorie Images) CreateCollection(collectionName string) (uint64, error) {
	statement, err := repositorie.db.Prepare("INSERT INTO collections(name) values(?)")
	if err != nil {
		return 0, fmt.Errorf("createcollection/db.prepare: %w", err)
	}
	defer statement.Close()

	result, err := statement.Exec(collectionName)
	if err != nil {
		return 0, fmt.Errorf("createcollection/statement.exec: %w", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createcollection/result.lastinsertedid: %w", err)
	}
	return uint64(lastInsertedId), nil
}

// CreateFace inserts a face value into the database.
func (repositorie Images) CreateFace(faceID string, imageID uint64) (uint64, error) {
	statement, err := repositorie.db.Prepare("INSERT INTO faces(face_id, image_id) values(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("createface/db.prepare: %w", err)
	}
	defer statement.Close()

	result, err := statement.Exec(faceID, imageID)
	if err != nil {
		return 0, fmt.Errorf("createface/statement.exec: %w", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createface/result.lastinsertedid: %w", err)
	}
	return uint64(lastInsertedId), nil
}

// CreateImage inserts a image value into the database.
func (repositorie Images) CreateImage(image models.Image, collectionID int) (uint64, error) {
	statement, err := repositorie.db.Prepare("INSERT INTO images(file_name, image_path, collection_id) values(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("createimage/db.prepare: %w", err)
	}
	defer statement.Close()

	result, err := statement.Exec(image.Filename, image.Path, collectionID)
	if err != nil {
		return 0, fmt.Errorf("createimage/statement.exec: %w", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createimage/result.lastinsertedid: %w", err)
	}
	return uint64(lastInsertedId), nil
}

// CreateMatch inserts a match between a face and a image into the database.
func (repositorie Images) CreateMatch(faceID, imageID uint64) error {
	statement, err := repositorie.db.Prepare(`INSERT INTO matches(face_id, image_id) values(?,?)`)
	if err != nil {
		return fmt.Errorf("creatematch/db.prepare: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(faceID, imageID)
	if err != nil {
		return fmt.Errorf("creatematch/statement.exec: %w", err)
	}

	return nil
}

// GetMatches queries the database for all the entries on the matches tables.
func (repositorie Images) GetMatches(faceID string) ([]models.Image, error) {
	lines, err := repositorie.db.Query(`
	SELECT i.file_name, i.image_path FROM matches m
	INNER JOIN faces f on f.id = m.face_id
	INNER JOIN images i on i.id = m.image_id
	WHERE f.face_id = ?; `, faceID)
	if err != nil {
		return nil, fmt.Errorf("getmatches/db.query: %w", err)
	}
	defer lines.Close()

	var images []models.Image
	for lines.Next() {
		var image models.Image
		if err = lines.Scan(
			&image.Filename,
			&image.Path,
		); err != nil {
			return nil, fmt.Errorf("getmatches/lines.scan: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

// CreateNoMatch inserts a image into the no_match table.
func (repositorie Images) CreateNoMatch(imageID uint64) error {
	statement, err := repositorie.db.Prepare(`INSERT INTO nomatches(image_id) values(?)`)
	if err != nil {
		return fmt.Errorf("createnomatch/db.prepare: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(imageID)
	if err != nil {
		return fmt.Errorf("createnomatch/statement.exec: %w", err)
	}

	return nil
}

// GetNoMatches returns all the entries on the no_matches table
func (repositorie Images) GetNoMatches() ([]models.Image, error) {
	lines, err := repositorie.db.Query(`
	SELECT i.file_name, i.image_path FROM nomatches nm 
	INNER JOIN images i ON  nm.image_id = i.id`)
	if err != nil {
		return nil, fmt.Errorf("getnomatches/db.query: %w", err)
	}
	defer lines.Close()

	var images []models.Image
	for lines.Next() {
		var image models.Image
		if err = lines.Scan(
			&image.Filename,
			&image.Path,
		); err != nil {
			return nil, fmt.Errorf("getnomatches/lines.scan: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

// GetFaceIDs returns all of the face_id entries from the faces table.
func (repositorie Images) GetFaceIDs(collectionID uint64) ([]models.Face, error) {
	lines, err := repositorie.db.Query(`
	SELECT face_id FROM faces f
	INNER JOIN images i ON f.image_id = i.id
	WHERE i.collection_id= ?`, collectionID)
	if err != nil {
		return nil, fmt.Errorf("getfaceids/db.query: %w", err)
	}
	defer lines.Close()

	var faces []models.Face
	for lines.Next() {
		var face models.Face
		if err := lines.Scan(
			&face.FaceID,
		); err != nil {
			return nil, fmt.Errorf("getfaceids/lines.scan: %w", err)
		}
		faces = append(faces, face)
	}

	return faces, nil
}

func (repositorie Images) GetCollectionID(collectionID string) (uint64, error) {
	line, err := repositorie.db.Query(`SELECT id FROM collections WHERE name = ?`, collectionID)
	if err != nil {
		return 0, fmt.Errorf("getcollectionid/db.query: %w", err)
	}
	defer line.Close()

	var ID int
	if line.Next() {
		if err := line.Scan(&ID); err != nil {
			return 0, fmt.Errorf("getcollectionid/lines.scan: %w", err)
		}
	}

	return uint64(ID), nil
}

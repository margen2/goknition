package repositories

import (
	"database/sql"
	"fmt"

	"github.com/margen2/goknition/models"
)

// Image represent a pointer to type sql.DB
type Images struct {
	db *sql.DB
}

// NewImagesRepositorie returns a pointer to type Images that allows database methods
func NewImagesRepositorie(db *sql.DB) *Images {
	return &Images{db}
}

// CreateFace inserts a face value into the database
func (repositorie Images) CreateFace(face models.Face, imageID uint64) error {
	statement, err := repositorie.db.Prepare("INSERT INTO faces(face_id, image_id) values(?, ?)")
	if err != nil {
		return fmt.Errorf("createface/db.prepare: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(face.ID, imageID)
	if err != nil {
		return fmt.Errorf("createface/statement.exec: %w", err)
	}

	return nil
}

// CreateImage inserts a image value into the database
func (repositorie Images) CreateImage(image models.Image) (uint64, error) {
	statement, err := repositorie.db.Prepare("INSERT INTO images(file_name, image_path) values (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("createimage/db.prepare: %w", err)
	}
	defer statement.Close()

	result, err := statement.Exec(image.FileName, image.Path)
	if err != nil {
		return 0, fmt.Errorf("createimage/statement.exec: %w", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createimage/result.lastinsertedid: %w", err)
	}
	return uint64(lastInsertedId), nil
}

// CreateMatch inserts a match between a face and a image into the database
func (repositorie Images) CreateMatch(faceID string, imageID string) error {
	statement, err := repositorie.db.Prepare(`
	INSERT INTO matches(face_id, image_id)
	SELECT (SELECT id FROM faces WHERE face_id = ? ),
	(SELECT id FROM images WHERE file_name = ?);`)
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

// GetMatches queries the database for all the entries on the matches tables
func (repositorie Images) GetMatches(faceID string) ([]models.Image, error) {
	lines, err := repositorie.db.Query(`
	SELECT i.file_name, i.image_path FROM matches m
	INNER JOIN faces f on f.id = m.face_id
	INNER JOIN images i on i.id = m.image_id
	WHERE f.face_id = ?;`, faceID)
	if err != nil {
		return nil, fmt.Errorf("getmatches/db.query: %w", err)
	}
	defer lines.Close()

	var images []models.Image
	for lines.Next() {
		var image models.Image
		if err = lines.Scan(
			&image.FileName,
			&image.Path,
		); err != nil {
			return nil, fmt.Errorf("getmatches/lines.scan: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

// CreateNoMatch inserts a image into the no_match table
func (repositorie Images) CreateNoMatch(fileName string) error {
	statement, err := repositorie.db.Prepare(`
	INSERT INTO nomatches(image_id)
	SELECT id FROM images WHERE file_name = ?;`)
	if err != nil {
		return fmt.Errorf("createnomatch/db.prepare: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(fileName)
	if err != nil {
		return fmt.Errorf("createnomatch/statement.exec: %w", err)
	}

	return nil
}

// GetNoMatches queries the database for all the entries on the no_matches table
func (repositorie Images) GetNoMatches() ([]models.Image, error) {
	lines, err := repositorie.db.Query(`SELECT i.file_name, i.image_path FROM nomatches nm 
	INNER JOIN images i ON  nm.image_id = i.id`)
	if err != nil {
		return nil, fmt.Errorf("getnomatches/db.query: %w", err)
	}
	defer lines.Close()

	var images []models.Image
	for lines.Next() {
		var image models.Image
		if err = lines.Scan(
			&image.FileName,
			&image.Path,
		); err != nil {
			return nil, fmt.Errorf("getnomatches/lines.scan: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

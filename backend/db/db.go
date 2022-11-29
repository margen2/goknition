package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var StringConnectionBD string

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", StringConnectionBD)
	if err != nil {
		return nil, fmt.Errorf("connectdb/sql.open: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connectdb/db.ping: %w", err)
	}

	return db, nil
}

func SetConnection(user, pw string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/?charset=utf8&parseTime=True&loc=Local", user, pw))
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS gok;")
	if err != nil {
		return err
	}
	db.Close()

	StringConnectionBD = fmt.Sprintf("%s:%s@/gok?charset=utf8&parseTime=True&loc=Local", user, pw)
	db, err = sql.Open("mysql", StringConnectionBD)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS collections(
		id INT auto_increment PRIMARY KEY,
		name varchar(50) UNIQUE NOT NULL
	) ENGINE=INNODB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS images (
		id INT auto_increment PRIMARY KEY,
		file_name VARCHAR(100) NOT NULL,
		file_path VARCHAR(200) NOT NULL,
		collection_id INT NOT NULL, 
		FOREIGN KEY (collection_id) 
		REFERENCES collections(id)
	) ENGINE=INNODB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS faces(
		id INT auto_increment PRIMARY KEY,
		face_id VARCHAR(100) NOT NULL, 
		collection_id INT NOT NULL 
		) ENGINE=INNODB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS matches(
		face_id INT NOT NULL,
		FOREIGN KEY (face_id)
		REFERENCES faces(id),
		image_id INT NOT NULL,
		FOREIGN KEY (image_id)
		REFERENCES images(id)
	) ENGINE=INNODB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS nomatches(
		image_id INT NOT NULL,
		FOREIGN KEY (image_id)
		REFERENCES images(id)
	) ENGINE=INNODB;`)
	if err != nil {
		return err
	}

	return nil
}

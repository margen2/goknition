CREATE DATABASE IF NOT EXISTS goknition;

USE goknition;

CREATE TABLE IF NOT EXISTS images (
    id INT auto_increment PRIMARY KEY,
    file_name VARCHAR(20) NOT NULL,
    image_path VARCHAR(150) NOT NULL
) ENGINE=INNODB;


CREATE TABLE IF NOT EXISTS faces(
    id INT auto_increment PRIMARY KEY,
    face_id VARCHAR(15) NOT NULL, 
    image_id INT NOT NULL, 
    FOREIGN KEY (image_id) 
    REFERENCES images(id)
) ENGINE=INNODB; 

CREATE TABLE IF NOT EXISTS matches(
    face_id INT NOT NULL,
    FOREIGN KEY (face_id)
    REFERENCES faces(id),
    image_id INT NOT NULL,
    FOREIGN KEY (image_id)
    REFERENCES images(id)
) ENGINE=INNODB; 

CREATE TABLE IF NOT EXISTS nomatches(
    image_id INT NOT NULL,
    FOREIGN KEY (image_id)
    REFERENCES images(id)
) ENGINE=INNODB; 


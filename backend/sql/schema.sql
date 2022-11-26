CREATE DATABASE IF NOT EXISTS goknition_new;

USE goknition_new;

CREATE TABLE IF NOT EXISTS collections(
    id INT auto_increment PRIMARY KEY,
    name varchar(50) UNIQUE NOT NULL
) ENGINE=INNODB; 


CREATE TABLE IF NOT EXISTS images (
    id INT auto_increment PRIMARY KEY,
    file_name VARCHAR(100) NOT NULL,
    file_path VARCHAR(200) NOT NULL,
    collection_id INT NOT NULL, 
    FOREIGN KEY (collection_id) 
    REFERENCES collections(id)
) ENGINE=INNODB;


CREATE TABLE IF NOT EXISTS faces(
    id INT auto_increment PRIMARY KEY,
    face_id VARCHAR(100) NOT NULL, 
    collection_id INT NOT NULL, 
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

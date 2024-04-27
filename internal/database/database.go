package database

import "database/sql"

func SetFile(path string, code string, db *sql.DB) error {
	insertCode := `
        INSERT INTO files (path, code) VALUES (?, ?);
	`

	_, err := db.Exec(insertCode, path, code)

	return err
}

func GetFile(code string, db *sql.DB) (string, error) {
	selectLink := `
	    SELECT path FROM files WHERE code = (?)
	`

	var link string
	err := db.QueryRow(selectLink, code).Scan(&link)

	return link, err
}

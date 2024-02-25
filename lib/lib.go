package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

func CreateFileAndUpdateDatabase(data []byte, db *sql.DB) (string, error) {
	path := fmt.Sprintf("./files/%v", uuid.New().String())

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	code := uuid.New().String()
	err = SetFile(path, code, db)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return code, nil
}

func RetrieveFile(code string, db *sql.DB) ([]byte, error) {
	path, err := GetFile(code, db)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)

	return data, err
}

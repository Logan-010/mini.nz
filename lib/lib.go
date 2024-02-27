package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

func PipelineErr[I any, O any](inputChan <-chan I, processer func(I) (O, error), errChan chan<- error) <-chan O {
	out := make(chan O)

	go func() {
		for i := range inputChan {
			output, err := processer(i)
			if err != nil {
				log.Println(err)
				errChan <- err
				break
			}

			out <- output
		}
		close(out)
	}()

	return out
}

func CreateFileAndUpdateDatabase(data []byte, db *sql.DB, filePath string, key string) (string, error) {
	errChan := make(chan error)
	encryptedDataChan := make(chan []byte)

	go func() {
		encryptedData, err := encrypt(data, []byte(key))
		if err != nil {
			errChan <- err
			return
		}

		encryptedDataChan <- encryptedData
	}()

	compressedDataChan := PipelineErr(encryptedDataChan, compress, errChan)

	path := fmt.Sprintf("%v/%v", filePath, uuid.New().String())

	go func() {
		file, err := os.Create(path)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()

		fileBytes := <-compressedDataChan

		_, err = file.Write(fileBytes)
		if err != nil {
			errChan <- err
			return
		}
	}()

	codeChan := make(chan string)

	go func() {
		code, err := uuid.NewRandom()
		if err != nil {
			errChan <- err
			return
		}

		uuidString := code.String()

		err = setFile(path, uuidString, db)
		if err != nil {
			errChan <- err
			return
		}

		codeChan <- uuidString
	}()

	select {
	case err := <-errChan:
		return "", err
	case code := <-codeChan:
		return code, nil
	}
}

func RetrieveFile(code string, db *sql.DB, key string) ([]byte, error) {
	errChan := make(chan error)
	pathChan := make(chan string)
	decryptedDataChan := make(chan []byte)

	go func() {
		path, err := getFile(code, db)
		if err != nil {
			errChan <- err
			return
		}

		pathChan <- path
	}()

	dataChan := PipelineErr(pathChan, os.ReadFile, errChan)
	uncompressedDataChan := PipelineErr(dataChan, decompress, errChan)

	go func() {
		decompressedData := <-uncompressedDataChan
		decryptedData, err := decrypt(decompressedData, []byte(key))
		if err != nil {
			errChan <- err
			return
		}

		decryptedDataChan <- decryptedData
	}()

	select {
	case err := <-errChan:
		return nil, err
	case finalData := <-decryptedDataChan:
		return finalData, nil
	}
}

package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"mini.nz/lib"
)

var db *sql.DB

func init() {
	err := os.Mkdir("./files", os.ModeDir)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS files (
        path TEXT NOT NULL,
        code TEXT NOT NULL
    );
    `

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
}

//go:embed assets/pages/404.html
var page404 string

//go:embed assets/pages/error.html
var pageInternalServiceError string

//go:embed assets/pages/index.html
var pageIndex string

//go:embed assets/style.css
var pageCss string

func errorCatcher(w http.ResponseWriter, _ *http.Request, code int) {
	w.WriteHeader(code)

	switch code {
	case http.StatusNotFound:
		fmt.Fprintln(w, page404)
	case http.StatusInternalServerError:
		fmt.Fprintln(w, pageInternalServiceError)

	default:
		fmt.Fprintln(w, pageInternalServiceError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorCatcher(w, r, http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, pageIndex)
}

func css(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/css" {
		errorCatcher(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/css")

	fmt.Fprintln(w, pageCss)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/upload" {
		errorCatcher(w, r, http.StatusNotFound)
		return
	}

	// limits to 50 mb memory
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	key := r.FormValue("encryptionKey")

	file, handle, err := r.FormFile("inputFile")
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// rejects file if it is over 100 mb
	if handle.Size > (100 << 20) {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	encryptedData, err := lib.Encrypt(fileBytes, []byte(key))
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	code, err := lib.CreateFileAndUpdateDatabase(encryptedData, db)
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, code)
}

func view(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	key := r.PathValue("key")

	data, err := lib.RetrieveFile(code, db)
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	decryptedData, err := lib.Decrypt(data, []byte(key))
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(decryptedData)
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/css", css)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/view/{code}/{key}", view)

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

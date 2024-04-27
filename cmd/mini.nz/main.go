package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"mini.nz/internal/compression"
	"mini.nz/internal/utils"
)

var (
	db          *sql.DB
	filesDirPtr *string = flag.String("filesDir", "./files", "Set flag to select the directory where user uploaded files are stored")
	dbDirPtr    *string = flag.String("dbDir", "./data.db", "Set flag to select the directory where the database is stored")
	portPtr     *uint   = flag.Uint("port", 8080, "Set flag to select port that app listens on")
	linkPtr     *string = flag.String("link", "http://localhost:8080", "Set to overide link inserted in html template. Please note that changing the link does not change the port the app listens on")
)

type PageTemplate struct {
	AppLink string
}

func initialize() {
	err := os.Mkdir(*filesDirPtr, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	db, err = sql.Open("sqlite3", *dbDirPtr)
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
		tmpl, err := template.New("404").Parse(page404)
		if err != nil {
			return
		}

		templateStruct := PageTemplate{*linkPtr}

		err = tmpl.Execute(w, templateStruct)
		if err != nil {
			return
		}
	case http.StatusInternalServerError:
		tmpl, err := template.New("internalServiceError").Parse(pageInternalServiceError)
		if err != nil {
			return
		}

		templateStruct := PageTemplate{*linkPtr}

		err = tmpl.Execute(w, templateStruct)
		if err != nil {
			return
		}

	default:
		tmpl, err := template.New("internalServiceError").Parse(pageInternalServiceError)
		if err != nil {
			return
		}

		templateStruct := PageTemplate{*linkPtr}

		err = tmpl.Execute(w, templateStruct)
		if err != nil {
			return
		}
	}
}

func css(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/css" {
		errorCatcher(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/css")

	fmt.Fprintln(w, pageCss)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorCatcher(w, r, http.StatusNotFound)
		return
	}

	tmpl, err := template.New("index").Parse(pageIndex)
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	templateStruct := PageTemplate{*linkPtr}

	err = tmpl.Execute(w, templateStruct)
	if err != nil {
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}
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

	code, err := utils.CreateFileAndUpdateDatabase(fileBytes, db, *filesDirPtr, key)
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

	data, err := utils.RetrieveFile(code, db, key)
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		errorCatcher(w, r, http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()
	initialize()

	addr := fmt.Sprintf(":%d", *portPtr)

	mux := http.NewServeMux()

	mux.Handle("/", compression.GzipHandler(http.HandlerFunc(index)))
	mux.Handle("/css", compression.GzipHandler(http.HandlerFunc(css)))
	mux.Handle("/upload", compression.GzipHandler(http.HandlerFunc(upload)))
	mux.Handle("/view/{code}/{key}", compression.GzipHandler(http.HandlerFunc(view)))

	fmt.Printf("Listening on %s (http://localhost%s)\n", *linkPtr, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

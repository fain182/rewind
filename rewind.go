package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fain182/rewind/archive"
	"github.com/goji/httpauth"
)

var records archive.RecordArchive

func main() {
	records = make(archive.RecordArchive)
	go webServer()
	for {
		go archive.Update(records)
		time.Sleep(30 * time.Minute)
	}
}

func webServer() {
	http.Handle("/", httpauth.SimpleBasicAuth(os.Getenv("USER"), os.Getenv("PASSWORD"))(http.HandlerFunc(handler)))
	println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(getAssetsPath() + "/index.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, archive.SortRecords(records))
}

func getAssetsPath() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath) + "/assets"
}

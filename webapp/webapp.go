package webapp

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fain182/rewind/archive"
	"github.com/goji/httpauth"
)

func Serve(records archive.RecordArchive) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { homepage(records, w, r) })
	handlerWithAuthentication := httpauth.SimpleBasicAuth(os.Getenv("USER"), os.Getenv("PASSWORD"))(handler)
	http.Handle("/", handlerWithAuthentication)
	println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(records archive.RecordArchive, w http.ResponseWriter, r *http.Request) {
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

package webapp

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fain182/rewind/archive"
	"github.com/goji/httpauth"
)

func Serve(records archive.RecordArchive) {
	homepageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { homepage(records, w, r) })
	homepageHandlerWithAuthentication := httpauth.SimpleBasicAuth(os.Getenv("USER"), os.Getenv("PASSWORD"))(homepageHandler)
	http.Handle("/", homepageHandlerWithAuthentication)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(records archive.RecordArchive, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(getAssetsPath() + "/index.html")
	if err != nil {
		panic(err)
	}

	jsonRecordings, err := json.Marshal(archive.SortRecords(records))
	if err != nil {
		log.Println(err)
	} else {
		t.Execute(w, template.JS(string(jsonRecordings)))
	}
}

func getAssetsPath() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath) + "/assets"
}

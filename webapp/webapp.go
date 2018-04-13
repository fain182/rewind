package webapp

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fain182/rewind/storage"
	"github.com/goji/httpauth"
)

// Serve start the webserver to show recordings data
func Serve(recordings storage.Recordings) {
	homepageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { homepage(recordings, w, r) })
	homepageHandlerWithAuthentication := httpauth.SimpleBasicAuth(os.Getenv("USER"), os.Getenv("PASSWORD"))(homepageHandler)
	http.Handle("/", homepageHandlerWithAuthentication)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(recordings storage.Recordings, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(getAssetsPath() + "/index.html")
	if err != nil {
		panic(err)
	}

	jsonRecordings, err := json.Marshal(recordings.Sort())
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

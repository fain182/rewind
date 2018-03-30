package main

import (
	"time"

	"github.com/fain182/rewind/scraper"
	"github.com/fain182/rewind/storage"
	"github.com/fain182/rewind/webapp"
)

var recordings storage.Recordings

func main() {
	recordings = make(storage.Recordings)
	go webapp.Serve(recordings)
	for {
		go scraper.Update(recordings)
		time.Sleep(30 * time.Minute)
	}
}

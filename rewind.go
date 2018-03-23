package main

import (
	"time"

	"github.com/fain182/rewind/archive"
	"github.com/fain182/rewind/webapp"
)

var records archive.RecordArchive

func main() {
	records = make(archive.RecordArchive)
	go webapp.Serve(records)
	for {
		go archive.Update(records)
		time.Sleep(30 * time.Minute)
	}
}

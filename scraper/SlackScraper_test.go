package scraper

import (
	"fmt"
	"testing"

	"github.com/fain182/rewind/storage"
	"github.com/nlopes/slack"
)

func TestAddOneMessageToRecordings(t *testing.T) {
	recordings := make(storage.Recordings)
	message := slack.File{
		Title:      "Registrazione incontro 06/03",
		Created:    slack.JSONTime(1508795665),
		Channels:   []string{"canale"},
		URLPrivate: "https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4",
	}

	addMessageToRecordings(recordings, message)

	if len(recordings) != 1 {
		t.Error("Expected 1 results, got ", len(recordings))
	}
	fmt.Println(recordings)
	r := recordings["1508795665000000000"]

	if r.Channel != "canale" {
		t.Error("Expected 'canale', got ", r.Channel)
	}

	if r.Title != "Registrazione incontro 06/03" {
		t.Error("Expected 'Registrazione incontro 06/03', got ", r.Title)
	}

	if r.URL != "https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4" {
		t.Error("Expected 'https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4', got ", r.URL)
	}
}

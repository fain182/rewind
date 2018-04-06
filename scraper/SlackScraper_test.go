package scraper

import (
	"testing"

	"github.com/fain182/rewind/storage"
	"github.com/nlopes/slack"
)

func TestAddOneMessageToRecordings(t *testing.T) {
	recordings := make(storage.Recordings)
	message := slack.SearchMessage{
		Text:      "uploaded a file: <https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4|Registrazione incontro 06/03>",
		Timestamp: "1508795665.000236",
		Channel:   slack.CtxChannel{ID: "123", Name: "Canale"},
	}

	addMessageToRecordings(recordings, message)

	if len(recordings) != 1 {
		t.Error("Expected 1 results, got ", len(recordings))
	}

	r := recordings[message.Timestamp]

	if r.Channel != "canale" {
		t.Error("Expected 'canale', got ", r.Channel)
	}

	if r.Title != "Registrazione incontro 06/03" {
		t.Error("Expected 'Registrazione incontro 06/03', got ", r.Title)
	}

	if r.Url != "https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4" {
		t.Error("Expected 'https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4', got ", r.Url)
	}
}

func TestTitleParsing(t *testing.T) {
	cases := []struct{ Body, ExpectedTitle string }{
		{
			"uploaded a file: <https://ideato.slack.com/files/U1KD1QEJ1/F9K5MT0PN/zoom_0.mp4|Registrazione incontro 06/03>",
			"Registrazione incontro 06/03",
		},
		{
			"uploaded a file: <https://ideato.slack.com/files/U024HSKKW/F8YHNGHCZ/zoom_0.mp4|zoom_0.mp4> and commented: @channel ecco la registrazione del lean coffee",
			"ecco la registrazione del lean coffee",
		},
		{
			"uploaded a file: <https://ideato.slack.com/files/U024HSKKW/F8YHNGHCZ/zoom_0.mp4|zoom_0.mp4>",
			"zoom 0",
		},
		{
			"uploaded a file: <https://ideato.slack.com/files/U024HSKKW/F8YHNGHCZ/common_language.m4a|common_language.m4a>",
			"common language",
		},
	}

	for _, c := range cases {
		actualTitle := parseTitle(c.Body)
		if actualTitle != c.ExpectedTitle {
			t.Error("Expected '", c.ExpectedTitle, "', got ", actualTitle)
		}
	}
}
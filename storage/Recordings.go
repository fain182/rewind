package storage

import (
	"math"
	"sort"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Recordings is the media collection
type Recordings map[string]Recording

// A Recording is a media uploaded to slack
type Recording struct {
	Title     string
	URL       string
	Channel   string
	Timestamp string
	HumanDate string
}

// Add a recording to the collection
func (recordings Recordings) Add(title, url, channel string, created time.Time) {
	recordings[strconv.FormatInt(created.UnixNano(), 10)] = Recording{
		Title:     title,
		URL:       url,
		Channel:   channel,
		Timestamp: strconv.FormatInt(created.Unix(), 10),
		HumanDate: humanize.Time(created),
	}
}

// Sort the recordings by creation date
func (recordings Recordings) Sort() []Recording {
	recordingList := make([]Recording, 0)
	for _, recording := range recordings {
		recordingList = append(recordingList, recording)
	}
	sort.Slice(recordingList[:], func(i, j int) bool {
		return recordingList[i].Timestamp > recordingList[j].Timestamp
	})
	return recordingList
}

func parseFloatTimestamp(timestamp string) string {
	timeFloat, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		panic(err)
	}
	sec, dec := math.Modf(timeFloat)
	uploadedAt := time.Unix(int64(sec), int64(dec*(1e9)))

	return humanize.Time(uploadedAt)
}

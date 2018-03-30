package storage

import (
	"math"
	"sort"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

type Recordings map[string]Recording

type Recording struct {
	Title     string
	Url       string
	Channel   string
	Timestamp string
	HumanDate string
}

func (recordings Recordings) Add(title, url, channel, timestamp string) {
	recordings[timestamp] = Recording{
		Title:     title,
		Url:       url,
		Channel:   channel,
		Timestamp: timestamp,
		HumanDate: parseFloatTimestamp(timestamp),
	}
}

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

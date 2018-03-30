package storage

import (
	"sort"
)

type Recordings map[string]Recording

type Recording struct {
	Title     string
	Url       string
	Channel   string
	Timestamp string
	HumanDate string
}

func SortRecordings(recordings Recordings) []Recording {
	recordingList := make([]Recording, 0)
	for _, recording := range recordings {
		recordingList = append(recordingList, recording)
	}
	sort.Slice(recordingList[:], func(i, j int) bool {
		return recordingList[i].Timestamp > recordingList[j].Timestamp
	})
	return recordingList
}

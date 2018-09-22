package scraper

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fain182/rewind/storage"
	"github.com/nlopes/slack"
)

// Update downloads recordings from slack and add them to storage
func Update(recordings storage.Recordings) {
	api := slack.New(os.Getenv("SLACK_API_KEY"))
	channels, err := api.GetChannels(true)
	if nil != err {
		log.Println(err.Error())
		return
	}

	searchFilesAndAddToRecordings(api, channels, recordings, "mp4")
	searchFilesAndAddToRecordings(api, channels, recordings, "m4a")
}

func searchFilesAndAddToRecordings(api *slack.Client, channels []slack.Channel, recordings storage.Recordings, query string) {
	files, paging, err := downloadMessageByPage(api, query, 0)
	addMessagesToRecordings(channels, recordings, files)
	if nil != err {
		log.Println(err.Error())
		return
	}

	for i := 1; i < paging.Count; i++ {
		time.Sleep(3 * time.Second)
		files, _, err := downloadMessageByPage(api, query, i)
		if nil != err {
			log.Println(err.Error())
			continue
		}
		addMessagesToRecordings(channels, recordings, files)
	}
}

func addMessagesToRecordings(channels []slack.Channel, recordings storage.Recordings, messages []slack.File) {
	for i := range messages {
		addMessageToRecordings(channels, recordings, messages[i])
	}
}

func addMessageToRecordings(channels []slack.Channel, recordings storage.Recordings, message slack.File) {
	channel := "[None]"
	if len(message.Channels) > 0 {
		for _, c := range channels {
			if c.ID == message.Channels[0] {
				channel = c.Name
				break
			}
		}
	}
	recordings.Add(normalizeTitle(message.Title), message.URLPrivate, channel, message.Created.Time())
}

func downloadMessageByPage(api *slack.Client, query string, pageNumber int) ([]slack.File, *slack.Paging, error) {
	params := slack.NewGetFilesParameters()
	params.Types = "videos,audios"
	params.Page = pageNumber
	return api.GetFiles(params)
}

func normalizeTitle(title string) string {
	title = strings.Replace(title, "_", " ", -1)
	title = strings.Replace(title, "-", " ", -1)
	return removeExtension(title)
}

func removeExtension(input string) string {
	r, _ := regexp.Compile("\\.\\w{3}$")
	return r.ReplaceAllString(input, "")
}

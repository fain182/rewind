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

	fileMessages, err := downloadMessageByPage(api, 0)
	addMessagesToRecordings(recordings, fileMessages.Matches)
	if nil != err {
		log.Println(err.Error())
		return
	}

	for i := 1; i < fileMessages.PageCount; i++ {
		time.Sleep(3 * time.Second)
		nextMessages, err := downloadMessageByPage(api, i)
		if nil != err {
			log.Println(err.Error())
			continue
		}
		addMessagesToRecordings(recordings, nextMessages.Matches)
	}

}

func addMessagesToRecordings(recordings storage.Recordings, messages []slack.File) {
	for i := range messages {
		addMessageToRecordings(recordings, messages[i])
	}
}

func addMessageToRecordings(recordings storage.Recordings, message slack.File) {
	channel := "[None]"
	if len(message.Channels) > 0 {
		channel = message.Channels[0]
	}
	recordings.Add(message.Title, message.URLPrivate, channel, message.Created.Time())
}

func downloadMessageByPage(api *slack.Client, pageNumber int) (*slack.SearchFiles, error) {
	params := slack.NewSearchParameters()
	params.Page = pageNumber
	return api.SearchFiles("mp4", params)
}

func isAZoomRecordingMessage(messageBody string) bool {
	return strings.Contains(messageBody, ".mp4") || strings.Contains(messageBody, ".m4a")
}

func parseTitle(messageBody string) string {
	beginTitleIndex := strings.Index(messageBody, "|") + 1
	endTitleIndex := strings.Index(messageBody, ">")
	defaultTitle := messageBody[beginTitleIndex:endTitleIndex]

	if defaultTitle == "zoom_0.mp4" {
		comment := strings.Trim(removeHandles(parseComment(messageBody)), " ")
		if comment != "" {
			return comment
		}
	}

	return normalizeTitle(defaultTitle)
}

func parseComment(messageBody string) string {
	commentPrefix := "and commented:"
	beginCommentIndex := strings.Index(messageBody, commentPrefix)

	if beginCommentIndex == -1 {
		return ""
	}

	return messageBody[beginCommentIndex+len(commentPrefix):]
}

func removeHandles(input string) string {
	r, _ := regexp.Compile("@\\w+")
	return r.ReplaceAllString(input, "")
}

func normalizeTitle(title string) string {
	title = strings.Replace(title, "_", " ", -1)
	return removeExtension(title)
}

func removeExtension(input string) string {
	r, _ := regexp.Compile("\\.\\w{3}$")
	return r.ReplaceAllString(input, "")
}

func parseURL(messageBody string) string {
	beginTitleIndex := strings.Index(messageBody, "<") + 1
	endTitleIndex := strings.LastIndex(messageBody, "|")
	return messageBody[beginTitleIndex:endTitleIndex]
}

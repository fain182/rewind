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

func Update(recordings storage.Recordings) {
	api := slack.New(os.Getenv("SLACK_API_KEY"))

	messages, err := downloadMessageByPage(api, 0)
	addMessagesToRecordings(recordings, messages.Matches)
	if nil != err {
		log.Println(err.Error())
		return
	}

	for i := 1; i < messages.PageCount; i++ {
		time.Sleep(3 * time.Second)
		nextMessages, err := downloadMessageByPage(api, i)
		if nil != err {
			log.Println(err.Error())
			continue
		}
		addMessagesToRecordings(recordings, nextMessages.Matches)
	}

}

func addMessagesToRecordings(recordings storage.Recordings, messages []slack.SearchMessage) {
	for i := range messages {
		addMessageToRecordings(recordings, messages[i])
	}
}

func addMessageToRecordings(recordings storage.Recordings, message slack.SearchMessage) {
	body := message.Text

	if isAZoomRecordingMessage(body) {
		recordings.Add(parseTitle(body), parseUrl(body), message.Channel.Name, message.Timestamp)
	}
}

func downloadMessageByPage(api *slack.Client, pageNumber int) (*slack.SearchMessages, error) {
	params := slack.NewSearchParameters()
	params.Page = pageNumber
	return api.SearchMessages("\"uploaded a file\"", params)
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

func parseUrl(messageBody string) string {
	beginTitleIndex := strings.Index(messageBody, "<") + 1
	endTitleIndex := strings.LastIndex(messageBody, "|")
	return messageBody[beginTitleIndex:endTitleIndex]
}

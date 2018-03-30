package scraper

import (
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
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
		recordings[message.Timestamp] = storage.Recording{
			Title:     getTitle(body),
			Url:       getUrl(body),
			Channel:   message.Channel.Name,
			Timestamp: message.Timestamp,
			HumanDate: getTime(message.Timestamp),
		}
	}
}

func downloadMessageByPage(api *slack.Client, pageNumber int) (*slack.SearchMessages, error) {
	params := slack.NewSearchParameters()
	params.Page = pageNumber
	return api.SearchMessages("\"uploaded a file\"", params)
}

func isAZoomRecordingMessage(messageBody string) bool {
	return strings.Contains(messageBody, "mp4") || strings.Contains(messageBody, "m4a")
}

func getTitle(messageBody string) string {
	beginTitleIndex := strings.Index(messageBody, "|") + 1
	endTitleIndex := strings.Index(messageBody, ">")
	defaultTitle := messageBody[beginTitleIndex:endTitleIndex]

	if defaultTitle == "zoom_0.mp4" {
		comment := strings.Trim(removeHandles(getComment(messageBody)), " ")
		if comment != "" {
			return comment
		}
	}

	return normalizeTitle(defaultTitle)
}

func getComment(messageBody string) string {
	commentPrefix := "and commented:"
	beginCommentIndex := strings.Index(messageBody, commentPrefix)

	if beginCommentIndex == -1 {
		return ""
	}

	return messageBody[beginCommentIndex+len(commentPrefix):]
}

func removeHandles(messageBody string) string {
	r, _ := regexp.Compile("@\\w+")
	return r.ReplaceAllString(messageBody, "")
}

func normalizeTitle(title string) string {
	title = strings.Replace(title, "_", " ", -1)
	return removeExtension(title)
}

func removeExtension(input string) string {
	r, _ := regexp.Compile("\\.\\w{3}$")
	return r.ReplaceAllString(input, "")
}

func getUrl(messageBody string) string {
	beginTitleIndex := strings.Index(messageBody, "<") + 1
	endTitleIndex := strings.LastIndex(messageBody, "|")
	return messageBody[beginTitleIndex:endTitleIndex]
}

func getTime(timestamp string) string {
	timeFloat, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		panic(err)
	}
	sec, dec := math.Modf(timeFloat)
	uploadedAt := time.Unix(int64(sec), int64(dec*(1e9)))

	return humanize.Time(uploadedAt)
}

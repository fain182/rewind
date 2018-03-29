package archive

import (
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/nlopes/slack"
)

type RecordArchive map[string]Record

type Record struct {
	Title     string
	Url       string
	Channel   string
	Timestamp string
	HumanDate string
}

func SortRecords(records RecordArchive) []Record {
	recordList := make([]Record, 0)
	for _, record := range records {
		recordList = append(recordList, record)
	}
	sort.Slice(recordList[:], func(i, j int) bool {
		return recordList[i].Timestamp > recordList[j].Timestamp
	})
	return recordList
}

func Update(records RecordArchive) {
	api := slack.New(os.Getenv("SLACK_API_KEY"))

	messages, err := downloadMessageByPage(api, 0)
	addMessagesToRecords(records, messages.Matches)
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
		addMessagesToRecords(records, nextMessages.Matches)
	}

}

func addMessagesToRecords(records RecordArchive, messages []slack.SearchMessage) {
	for i := range messages {
		addMessageToRecords(records, messages[i])
	}
}

func addMessageToRecords(records RecordArchive, message slack.SearchMessage) {
	body := message.Text

	if isAZoomRecordingMessage(body) {
		records[message.Timestamp] = Record{
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

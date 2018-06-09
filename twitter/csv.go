package twitter

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"
)

type CSVTweet struct {
	TweetID                  int64
	InReplyToStatusID        int64
	InReplyToUserID          int64
	Timestamp                time.Time
	Source                   string
	Text                     string
	RetweetedStatusID        int64
	RetweetedStatusUserID    int64
	RetweetedStatusTimestamp time.Time
	ExpandedURLs             []string
}

type Reader struct {
	r       *csv.Reader
	skipped bool
}

func parseInt(i string) (int64, error) {
	if i == "" {
		return 0, nil
	}
	return strconv.ParseInt(i, 10, 64)
}

func NewCSVReader(in io.Reader) *Reader {
	return &Reader{r: csv.NewReader(in)}
}

func (c *Reader) Read() (r CSVTweet, err error) {
	record, err := c.r.Read()
	if !c.skipped {
		record, err = c.r.Read()
		c.skipped = true
	}

	if err != nil {
		return CSVTweet{}, err
	}

	r.TweetID, err = parseInt(record[0])
	if err != nil {
		return
	}

	r.InReplyToStatusID, err = parseInt(record[1])
	if err != nil {
		return
	}

	r.InReplyToUserID, err = parseInt(record[2])
	if err != nil {
		return
	}

	if record[3] != "" {
		r.Timestamp, err = time.Parse("2006-01-02 15:04:05 -0700", record[3])
		if err != nil {
			return
		}
	}

	r.RetweetedStatusID, err = parseInt(record[6])
	if err != nil {
		return
	}

	r.RetweetedStatusUserID, err = parseInt(record[7])
	if err != nil {
		return
	}

	if record[8] != "" {
		r.RetweetedStatusTimestamp, err = time.Parse("2006-01-02 15:04:05 -0700", record[8])
		if err != nil {
			return
		}
	}

	r.Source = record[4]
	r.Text = record[5]
	r.ExpandedURLs = strings.Split(record[0], ",")

	return
}

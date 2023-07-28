package rss

import (
	"encoding/xml"
	"io"
	"module36a/pkg/storage"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type rssFeed struct {
	Rss     string  `xml:"rss"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"description"`
	Link    string `xml:"link"`
	PubTime string `xml:"pubDate"`
}

// ParseRss читает новости из rss рассылки url с интервалом period
// и отправляет их в chan posts
func ParseRss(url string, db storage.Interfase, period int, posts chan<- []storage.Post, errs chan<- error) {
	for {
		news, err := Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}

// Parse возвращает слайс новостей из rss рассылки из url
func Parse(url string) ([]storage.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed rssFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	var news []storage.Post
	p := bluemonday.NewPolicy()

	for _, item := range feed.Channel.Items {
		var post storage.Post
		post.Title = item.Title
		item.Content = p.Sanitize(item.Content)
		post.Content = item.Content
		post.Link = item.Link

		// t, err := time.Parse(time.RFC1123Z, item.PubTime) // RFC1123Z=="Mon, 02 Jan 2006 15:04:05 -0700"
		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubTime)
		if err == nil {
			post.PubTime = t.Unix()
		}
		t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.PubTime) // RFC1123=="Mon, 02 Jan 2006 15:04:05 MST"
		if err == nil {
			post.PubTime = t.Unix()
		}
		news = append(news, post)
	}

	return news, nil
}

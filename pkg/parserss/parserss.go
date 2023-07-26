package parserss

import (
	"module36a/pkg/rss"
	"module36a/pkg/storage"
	"time"
)

// Parse читает новости из rss рассылки url с интервалом period
// и отправляет их в chan posts
func Parse(url string, db storage.Interfase, period int, posts chan<- []storage.Post, errs chan<- error) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}

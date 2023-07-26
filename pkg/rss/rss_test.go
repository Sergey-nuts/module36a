package rss

import (
	"testing"
)

func TestParse(t *testing.T) {
	url := "https://habr.com/ru/rss/hub/go/all/?fl=ru"
	news, err := Parse(url)
	if err != nil {
		t.Fatal(err)
	}
	if len(news) == 0 {
		t.Fatal("data not unmarshal")
	}
	t.Logf("received %v news: %v", len(news), news)
}

package memdb

import (
	"module36a/pkg/storage"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db := New()

	news := []storage.Post{
		{Title: "first post", Content: "first content", PubTime: time.Now().Unix(), Link: "http://test.url.com/post1"},
		{Title: "second post", Content: "second content", PubTime: time.Now().Unix(), Link: "http://testing.url.com/post2"},
	}

	err := db.AddNews(news)
	if err != nil {
		t.Fatal(err)
	}

	got, err := db.News(2)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("memdb.news() got=%v, want=%v", got, news)
	}

}

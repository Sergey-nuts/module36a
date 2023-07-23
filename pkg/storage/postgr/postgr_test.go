package postgr

import (
	"module36a/pkg/storage"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	// for example docker cli:
	// docker run --rm -it -p 5432:5432/tcp --name postgres -e POSTGRES_PASSWORD=testpasswd postgres

	db, err := New("postgres://postgres:testpasswd@localhost/testDB")
	// postgrUser := os.Getenv("dbuser")
	// postgrPwd := os.Getenv("dbpass")
	// dbhost := os.Getenv("dbhost")
	// db, err := New("postgres://" + postgrUser + ":" + postgrPwd + "@" + dbhost + "/module36a-test")
	if err != nil {
		t.Fatal(err)
	}

	news := []storage.Post{
		{Title: "first post", Content: "first content", PubTime: time.Now().Unix(), Link: "http://test.url.com/post1"},
		{Title: "second post", Content: "second content", PubTime: time.Now().Unix(), Link: "http://testing.url.com/post2"},
	}

	err = db.AddNews(news)
	if err != nil {
		t.Fatal(err)
	}

	got, err := db.News(2)
	if err != nil {
		t.Fatal(err)
	}
	news[0].ID = got[0].ID
	news[1].ID = got[1].ID
	if news[0] != got[1] && news[1] != got[0] /*!reflect.DeepEqual(news, got)*/ {
		t.Fatalf("postgr.news got=%v, want=%v", got, news)
	}
}

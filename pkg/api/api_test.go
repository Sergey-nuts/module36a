package api

import (
	"encoding/json"
	"io"
	"module36a/pkg/storage"
	"module36a/pkg/storage/memdb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_postsHandler(t *testing.T) {
	db := memdb.New()
	news := []storage.Post{
		{Title: "first post", Content: "first content", PubTime: time.Now().Unix(), Link: "http://test.url.com/post1"},
		{Title: "second post", Content: "second content", PubTime: time.Now().Unix(), Link: "http://testing.url.com/post2"},
	}
	db.AddNews(news)

	api := New(db)

	req := httptest.NewRequest(http.MethodGet, "/news/2", nil)
	w := httptest.NewRecorder()

	api.r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Response status code: %v, want: %v", w.Code, http.StatusOK)
	}

	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Unnable to read body: %v", err)
	}

	var got []storage.Post
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("unnable to unmarshal body: %v", err)
	}

	news[0].ID = 1
	news[1].ID = 2
	if news[0] != got[1] && news[1] != got[0] {
		t.Errorf("got: %v, want: %v", got, news)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"module36a/pkg/api"
	"module36a/pkg/rss"
	"module36a/pkg/storage"
	"module36a/pkg/storage/postgr"
)

type config struct {
	RssUrls []string `json:"rss"`
	Period  int      `json:"request_period"`
}

var confFile = "./config.json"

func main() {
	// db := memdb.New()

	postgrUser := os.Getenv("dbuser")
	postgrPwd := os.Getenv("dbpass")
	dbhost := os.Getenv("dbhost")
	log.Printf("conncting to postgresql... ")
	db, err := postgr.New("postgres://" + postgrUser + ":" + postgrPwd + "@" + dbhost + "/module36a")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("postgresql connected.\n")
	defer db.Close()

	api := api.New(db)

	// чтение файла с настройками
	b, err := os.ReadFile(confFile)
	if err != nil {
		log.Fatal(err)
	}
	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error)           // канал для ошибок
	posts := make(chan []storage.Post) // канал для новостей

	// чтение rss каналов в отдельных горутинах
	for _, url := range conf.RssUrls {
		go rss.ParseRss(url, db, conf.Period, posts, errs)
	}

	// логирование ошибок
	go func() {
		for err := range errs {
			log.Println("gonews error:", err)
		}
	}()

	// добавление новостей в базу данных
	go func() {
		for news := range posts {
			err := db.AddNews(news)
			if err != nil {
				errs <- err
			}
		}
	}()

	// запуск веб-сервера с API и приложением
	log.Println("starting server")
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

package memdb

import (
	"fmt"
	"module36a/pkg/storage"
	"sync"
)

type DB struct {
	m    sync.Mutex
	id   int                  // текущее значение ID для новой записи
	news map[int]storage.Post // БД
}

// конструктор БД
func New() *DB {
	db := DB{
		id:   1,
		news: map[int]storage.Post{},
	}

	return &db
}

// AddNews добавляет новости из news в базу даннх
func (d *DB) AddNews(news []storage.Post) error {
	for _, post := range news {
		post.ID = d.id
		d.add(post)
		d.id++
	}

	return nil
}

// add добавляет p в базу данных
func (d *DB) add(p storage.Post) {
	d.m.Lock()
	defer d.m.Unlock()
	d.news[p.ID] = p
}

// News возвращает последние n новостенй из базы данных
func (d *DB) News(n int) ([]storage.Post, error) {
	if n >= d.id {
		return nil, fmt.Errorf("unable to get news, n=%v bigger then news id=%v", n, d.id)
	}
	d.m.Lock()
	defer d.m.Unlock()
	var news []storage.Post
	for i := n; i > 0; i-- {
		news = append(news, d.news[i])
	}

	return news, nil
}

package postgr

import (
	"context"
	"fmt"
	"module36a/pkg/storage"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

// конструктор БД
func New(conf string) (*Postgres, error) {
	db, err := pgxpool.Connect(context.Background(), conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{db}, nil
}

// AddNews добавляет новости из news в базу даннх
func (p *Postgres) AddNews(news []storage.Post) error {
	ctx := context.Background()
	for _, post := range news {
		_, err := p.db.Exec(ctx, `
			INSERT INTO news(title, content, pubtime, link)
			VALUES ($1, $2, $3, $4);
		`, post.Title, post.Content, post.PubTime, post.Link)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// News возвращает последние n новостенй из базы данных
func (p *Postgres) News(n int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), `
		SELECT id, title, content, pubtime, link 
		FROM news
		ORDER BY pubtime DESC
		LIMIT $1;
	`,
		n,
	)
	if err != nil {
		return nil, err
	}

	var news []storage.Post
	var post storage.Post
	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}
		news = append(news, post)
	}

	return news, rows.Err()
}

// Close закрывает все соединения с базой
func (p *Postgres) Close() {
	p.db.Close()
}

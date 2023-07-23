package postgr

import (
	"context"
	"module36a/pkg/storage"

	"github.com/jackc/pgx/v4"
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
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// пакетный запрос
	batch := new(pgx.Batch)
	for _, post := range news {
		batch.Queue(`
			INSERT INTO news(title, content, pubtime, link)
			VALUES ($1, $2, $3, $4);
		`,
			post.Title, post.Content, post.PubTime, post.Link,
		)
	}
	res := tx.SendBatch(ctx, batch)

	err = res.Close()
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// News возвращает последние n новостенй из базы данных
func (p *Postgres) News(n int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), `
		SELECT id, title, content, pubtime, link 
		FROM news
		ORDER BY pubtime
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

package storage

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

type Interfase interface {
	AddNews([]Post) error
	News(int) ([]Post, error)
}

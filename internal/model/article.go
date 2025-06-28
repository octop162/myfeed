package model

import "time"

// Article は記事のデータモデルです。
type Article struct {
	ID         string    `json:"id"`
	FeedID     string    `json:"feed_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content,omitempty"`
	URL        string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	IsRead     bool      `json:"is_read"`
	IsLater    bool      `json:"is_later"`
	CreatedAt  time.Time `json:"created_at"`
}

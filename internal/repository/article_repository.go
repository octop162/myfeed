package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"feedapp/internal/model"
)

// ArticleRepository は記事のデータ永続化を定義するインターフェースです。
type ArticleRepository interface {
	GetAll() ([]model.Article, error)
	GetByID(id string) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Update(article model.Article) (model.Article, error)
	Delete(id string) error
	GetLaterArticles() ([]model.Article, error)
}

// articleRepository は ArticleRepository インターフェースの実装です。
type articleRepository struct {
	db *sql.DB
}

// NewArticleRepository は新しい articleRepository インスタンスを作成します。
func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetAll() ([]model.Article, error) {
	rows, err := r.db.Query("SELECT id, feed_id, title, content, url, published_at, is_read, is_later, created_at FROM articles")
	if err != nil {
		return nil, fmt.Errorf("failed to get all articles: %w", err)
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var content sql.NullString
		var publishedAt sql.NullTime
		if err := rows.Scan(&article.ID, &article.FeedID, &article.Title, &content, &article.URL, &publishedAt, &article.IsRead, &article.IsLater, &article.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		if content.Valid {
			article.Content = content.String
		}
		if publishedAt.Valid {
			article.PublishedAt = publishedAt.Time
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return articles, nil
}

func (r *articleRepository) GetByID(id string) (model.Article, error) {
	var article model.Article
	var content sql.NullString
	var publishedAt sql.NullTime
	err := r.db.QueryRow("SELECT id, feed_id, title, content, url, published_at, is_read, is_later, created_at FROM articles WHERE id = $1", id).Scan(&article.ID, &article.FeedID, &article.Title, &content, &article.URL, &publishedAt, &article.IsRead, &article.IsLater, &article.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Article{}, ErrNotFound
		}
		return model.Article{}, fmt.Errorf("failed to get article by ID: %w", err)
	}
	if content.Valid {
		article.Content = content.String
	}
	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}
	return article, nil
}

func (r *articleRepository) Create(article model.Article) (model.Article, error) {
	query := `INSERT INTO articles (id, feed_id, title, content, url, published_at, is_read, is_later, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, feed_id, title, content, url, published_at, is_read, is_later, created_at`
	var createdArticle model.Article
	var content sql.NullString
	var publishedAt sql.NullTime
	err := r.db.QueryRow(query, article.ID, article.FeedID, article.Title, sql.NullString{String: article.Content, Valid: article.Content != ""}, article.URL, sql.NullTime{Time: article.PublishedAt, Valid: !article.PublishedAt.IsZero()}, article.IsRead, article.IsLater, article.CreatedAt).Scan(&createdArticle.ID, &createdArticle.FeedID, &createdArticle.Title, &content, &createdArticle.URL, &publishedAt, &createdArticle.IsRead, &createdArticle.IsLater, &createdArticle.CreatedAt)
	if err != nil {
		return model.Article{}, fmt.Errorf("failed to create article: %w", err)
	}
	if content.Valid {
		createdArticle.Content = content.String
	}
	if publishedAt.Valid {
		createdArticle.PublishedAt = publishedAt.Time
	}
	return createdArticle, nil
}

func (r *articleRepository) Update(article model.Article) (model.Article, error) {
	query := `UPDATE articles SET title = $1, content = $2, url = $3, published_at = $4, is_read = $5, is_later = $6 WHERE id = $7 RETURNING id, feed_id, title, content, url, published_at, is_read, is_later, created_at`
	var updatedArticle model.Article
	var content sql.NullString
	var publishedAt sql.NullTime
	err := r.db.QueryRow(query, article.Title, sql.NullString{String: article.Content, Valid: article.Content != ""}, article.URL, sql.NullTime{Time: article.PublishedAt, Valid: !article.PublishedAt.IsZero()}, article.IsRead, article.IsLater, article.ID).Scan(&updatedArticle.ID, &updatedArticle.FeedID, &updatedArticle.Title, &content, &updatedArticle.URL, &publishedAt, &updatedArticle.IsRead, &updatedArticle.IsLater, &updatedArticle.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Article{}, ErrNotFound
		}
		return model.Article{}, fmt.Errorf("failed to update article: %w", err)
	}
	if content.Valid {
		updatedArticle.Content = content.String
	}
	if publishedAt.Valid {
		updatedArticle.PublishedAt = publishedAt.Time
	}
	return updatedArticle, nil
}

func (r *articleRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepository) GetLaterArticles() ([]model.Article, error) {
	rows, err := r.db.Query("SELECT id, feed_id, title, content, url, published_at, is_read, is_later, created_at FROM articles WHERE is_later = TRUE")
	if err != nil {
		return nil, fmt.Errorf("failed to get later articles: %w", err)
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var content sql.NullString
		var publishedAt sql.NullTime
		if err := rows.Scan(&article.ID, &article.FeedID, &article.Title, &content, &article.URL, &publishedAt, &article.IsRead, &article.IsLater, &article.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		if content.Valid {
			article.Content = content.String
		}
		if publishedAt.Valid {
			article.PublishedAt = publishedAt.Time
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return articles, nil
}

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"feedapp/internal/model"
)

// FeedRepository はフィードのデータ永続化を定義するインターフェースです。
type FeedRepository interface {
	GetAll() ([]model.Feed, error)
	GetByID(id string) (model.Feed, error)
	Create(feed model.Feed) (model.Feed, error)
	Update(feed model.Feed) (model.Feed, error)
	Delete(id string) error
}

// feedRepository は FeedRepository インターフェースの実装です。
type feedRepository struct {
	db *sql.DB
}

// NewFeedRepository は新しい feedRepository インスタンスを作成します。
func NewFeedRepository(db *sql.DB) FeedRepository {
	return &feedRepository{db: db}
}

func (r *feedRepository) GetAll() ([]model.Feed, error) {
	rows, err := r.db.Query("SELECT id, name, url, plugin_type, folder_id, update_interval, last_updated, created_at FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("failed to get all feeds: %w", err)
	}
	defer rows.Close()

	var feeds []model.Feed
	for rows.Next() {
		var feed model.Feed
		var folderID sql.NullString
		var lastUpdated sql.NullTime
		if err := rows.Scan(&feed.ID, &feed.Name, &feed.URL, &feed.PluginType, &folderID, &feed.UpdateInterval, &lastUpdated, &feed.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan feed row: %w", err)
		}
		if folderID.Valid {
			feed.FolderID = folderID.String
		}
		if lastUpdated.Valid {
			feed.LastUpdated = lastUpdated.Time
		}
		feeds = append(feeds, feed)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return feeds, nil
}

func (r *feedRepository) GetByID(id string) (model.Feed, error) {
	var feed model.Feed
	var folderID sql.NullString
	var lastUpdated sql.NullTime
	err := r.db.QueryRow("SELECT id, name, url, plugin_type, folder_id, update_interval, last_updated, created_at FROM feeds WHERE id = $1", id).Scan(&feed.ID, &feed.Name, &feed.URL, &feed.PluginType, &folderID, &feed.UpdateInterval, &lastUpdated, &feed.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Feed{}, ErrNotFound
		}
		return model.Feed{}, fmt.Errorf("failed to get feed by ID: %w", err)
	}
	if folderID.Valid {
		feed.FolderID = folderID.String
	}
	if lastUpdated.Valid {
		feed.LastUpdated = lastUpdated.Time
	}
	return feed, nil
}

func (r *feedRepository) Create(feed model.Feed) (model.Feed, error) {
		query := `INSERT INTO feeds (id, name, url, plugin_type, folder_id, update_interval, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, url, plugin_type, folder_id, update_interval, last_updated, created_at`
	var createdFeed model.Feed
	var folderID sql.NullString
	var lastUpdated sql.NullTime
	err := r.db.QueryRow(query, feed.ID, feed.Name, feed.URL, feed.PluginType, sql.NullString{String: feed.FolderID, Valid: feed.FolderID != ""}, feed.UpdateInterval, feed.CreatedAt).Scan(&createdFeed.ID, &createdFeed.Name, &createdFeed.URL, &createdFeed.PluginType, &folderID, &createdFeed.UpdateInterval, &lastUpdated, &createdFeed.CreatedAt)
	if err != nil {
		return model.Feed{}, fmt.Errorf("failed to create feed: %w", err)
	}
	if folderID.Valid {
		createdFeed.FolderID = folderID.String
	}
	if lastUpdated.Valid {
		createdFeed.LastUpdated = lastUpdated.Time
	}
	return createdFeed, nil
}

func (r *feedRepository) Update(feed model.Feed) (model.Feed, error) {
	query := `UPDATE feeds SET name = $1, url = $2, plugin_type = $3, folder_id = $4, update_interval = $5, last_updated = $6 WHERE id = $7 RETURNING id, name, url, plugin_type, folder_id, update_interval, last_updated, created_at`
	var updatedFeed model.Feed
	var folderID string
	var lastUpdated time.Time
	err := r.db.QueryRow(query, updatedFeed.Name, updatedFeed.URL, updatedFeed.PluginType, updatedFeed.FolderID, updatedFeed.UpdateInterval, updatedFeed.LastUpdated, updatedFeed.ID).Scan(&updatedFeed.ID, &updatedFeed.Name, &updatedFeed.URL, &updatedFeed.PluginType, &folderID, &updatedFeed.UpdateInterval, &lastUpdated, &updatedFeed.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Feed{}, ErrNotFound
		}
		return model.Feed{}, fmt.Errorf("failed to update feed: %w", err)
	}
	updatedFeed.FolderID = folderID
	updatedFeed.LastUpdated = lastUpdated
	return updatedFeed, nil
}

func (r *feedRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM feeds WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete feed: %w", err)
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

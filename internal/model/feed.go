package model

import "time"

// Feed はフィードのデータモデルです。
type Feed struct {
	ID             string    `json:"id"`
	Name           string    `json:"name" binding:"required"`
	URL            string    `json:"url" binding:"required,url"`
	PluginType     string    `json:"plugin_type" binding:"required"`
	FolderID       string    `json:"folder_id,omitempty"`
	UpdateInterval int       `json:"update_interval"`
	LastUpdated    time.Time `json:"last_updated,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

package model

import (
	"time"

	"github.com/google/uuid" // uuidパッケージを追加
)

// Folder はフォルダのデータモデルです。
type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	UserID    string    `json:"user_id,omitempty"` // 将来対応
	CreatedAt time.Time `json:"created_at"`
}

// GenerateUUID は新しいUUIDを生成して文字列で返します。
func GenerateUUID() string {
	return uuid.New().String()
}


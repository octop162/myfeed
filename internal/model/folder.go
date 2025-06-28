// Package model はフィードアプリケーションのデータモデルを定義します。
//
// このパッケージには、データベースで使用される構造体とそれに関連する
// ユーティリティ関数が含まれています。
package model

import (
	"time"

	"github.com/google/uuid"
)

// Folder はフォルダのデータモデルです。
//
// フォルダはフィードを分類するためのコンテナとして機能し、
// ユーザーがフィードを整理するために使用されます。
// 各フォルダは一意のIDを持ち、複数のフィードを含むことができます。
//
// JSON tags:
//   - id: フォルダの一意識別子（UUID形式）
//   - name: フォルダの名前（必須フィールド）
//   - user_id: 所有者のユーザーID（将来のマルチユーザー対応用、現在は未使用）
//   - created_at: フォルダの作成日時
type Folder struct {
	ID        string    `json:"id"`                   // フォルダの一意識別子
	Name      string    `json:"name" binding:"required"` // フォルダ名（必須）
	UserID    string    `json:"user_id,omitempty"`    // 将来のマルチユーザー対応用
	CreatedAt time.Time `json:"created_at"`           // 作成日時
}

// GenerateUUID は新しいUUIDを生成して文字列で返します。
//
// この関数はGoogle UUID パッケージを使用してランダムなUUID（v4）を生成し、
// 文字列形式で返します。主にデータベースレコードの一意識別子として
// 使用されます。
//
// Returns:
//   - string: 新しく生成されたUUID文字列（例: "550e8400-e29b-41d4-a716-446655440000"）
func GenerateUUID() string {
	return uuid.New().String()
}


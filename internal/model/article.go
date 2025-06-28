package model

import "time"

// Article は記事のデータモデルです。
//
// 記事はフィードから取得される個々のコンテンツを表現します。
// 各記事は読了状態と「後で見る」状態を管理し、ユーザーの読書習慣を
// サポートします。記事は削除されることなく永続的に保存されます。
//
// 状態管理:
//   - 未読 + 通常
//   - 未読 + 後で見る
//   - 既読 + 通常
//   - 既読 + 後で見る
//
// JSON tags:
//   - id: 記事の一意識別子（UUID形式）
//   - feed_id: 記事が属するフィードのID
//   - title: 記事のタイトル
//   - content: 記事の本文（省略可能）
//   - url: 記事の元URLリンク
//   - published_at: 記事の公開日時
//   - is_read: 既読フラグ
//   - is_later: 後で見るフラグ
//   - created_at: 記事の作成日時
type Article struct {
	ID         string    `json:"id"`                   // 記事の一意識別子
	FeedID     string    `json:"feed_id"`              // 所属フィードID
	Title      string    `json:"title"`                // 記事タイトル
	Content    string    `json:"content,omitempty"`    // 記事本文（省略可能）
	URL        string    `json:"url"`                  // 記事の元URL
	PublishedAt time.Time `json:"published_at"`        // 公開日時
	IsRead     bool      `json:"is_read"`              // 既読フラグ
	IsLater    bool      `json:"is_later"`             // 後で見るフラグ
	CreatedAt  time.Time `json:"created_at"`           // 作成日時
}

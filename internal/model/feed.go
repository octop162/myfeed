package model

import "time"

// Feed はフィードのデータモデルです。
//
// フィードはRSSやカスタムスクレイピングによる情報源を表現します。
// 各フィードは特定のURLから定期的に記事を取得し、指定されたフォルダに分類されます。
// プラグインシステムにより、RSS以外の様々なソースにも対応可能です。
//
// JSON tags:
//   - id: フィードの一意識別子（UUID形式）
//   - name: フィードの表示名（必須フィールド）
//   - url: フィードのURL（必須、URL形式の検証あり）
//   - plugin_type: 使用するプラグインの種別（例: "rss", "custom"）
//   - folder_id: 所属するフォルダのID（省略可能）
//   - update_interval: 更新間隔（分単位）
//   - last_updated: 最後に更新された日時
//   - created_at: フィードの作成日時
type Feed struct {
	ID             string    `json:"id"`                        // フィードの一意識別子
	Name           string    `json:"name" binding:"required"`   // フィード名（必須）
	URL            string    `json:"url" binding:"required,url"` // フィードURL（必須、URL形式）
	PluginType     string    `json:"plugin_type" binding:"required"` // プラグイン種別（必須）
	FolderID       string    `json:"folder_id,omitempty"`       // 所属フォルダID
	UpdateInterval int       `json:"update_interval"`           // 更新間隔（分）
	LastUpdated    time.Time `json:"last_updated,omitempty"`    // 最終更新日時
	CreatedAt      time.Time `json:"created_at"`                // 作成日時
}

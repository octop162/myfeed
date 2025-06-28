# テーブル設計

## 概要
カスタムフィードアプリケーションのデータベーステーブル設計を定義する。

## テーブル一覧

### `folders` テーブル
- **説明**: フィードを分類するためのフォルダ情報を格納する。
- **カラム**:
    - `id` (UUID): プライマリキー。自動生成。
    - `name` (VARCHAR(255)): フォルダ名。NULL不可。
    - `user_id` (UUID): ユーザーID。将来的なマルチユーザー対応のための外部キー（現在はNULL許容）。
    - `created_at` (TIMESTAMP WITH TIME ZONE): レコード作成日時。デフォルトは現在時刻。

### `feeds` テーブル
- **説明**: 登録されたフィードソースの情報を格納する。
- **カラム**:
    - `id` (UUID): プライマリキー。自動生成。
    - `name` (VARCHAR(255)): フィード名。NULL不可。
    - `url` (TEXT): フィードのURL。NULL不可、ユニーク。
    - `plugin_type` (VARCHAR(255)): 使用するプラグインの種別（例: 'rss', 'custom'）。NULL不可。
    - `folder_id` (UUID): 所属するフォルダのID。`folders`テーブルの`id`を参照。フォルダが削除された場合はNULLになる。
    - `update_interval` (INTEGER): 更新間隔（分）。デフォルトは360分（6時間）。
    - `last_updated` (TIMESTAMP WITH TIME ZONE): 最終更新日時。
    - `created_at` (TIMESTAMP WITH TIME ZONE): レコード作成日時。デフォルトは現在時刻。

### `articles` テーブル
- **説明**: 各フィードから取得された記事の情報を格納する。
- **カラム**:
    - `id` (UUID): プライマリキー。自動生成。
    - `feed_id` (UUID): 所属するフィードのID。`feeds`テーブルの`id`を参照。フィードが削除された場合は記事も削除される。
    - `title` (TEXT): 記事のタイトル。NULL不可。
    - `content` (TEXT): 記事の本文。NULL許容。
    - `url` (TEXT): 記事のオリジナルURL。NULL不可、ユニーク。
    - `published_at` (TIMESTAMP WITH TIME ZONE): 記事の公開日時。
    - `is_read` (BOOLEAN): 記事が既読かどうか。デフォルトはFALSE。
    - `is_later` (BOOLEAN): 記事が「後で見る」に設定されているか。デフォルトはFALSE。
    - `created_at` (TIMESTAMP WITH TIME ZONE): レコード作成日時。デフォルトは現在時刻。

### `plugins` テーブル
- **説明**: カスタムスクレイピングプラグインの情報を格納する。
- **カラム**:
    - `id` (UUID): プライマリキー。自動生成。
    - `name` (VARCHAR(255)): プラグイン名。NULL不可、ユニーク。
    - `file_path` (TEXT): プラグインの実行可能ファイルパス。NULL不可。
    - `enabled` (BOOLEAN): プラグインが有効かどうか。デフォルトはTRUE。
    - `created_at` (TIMESTAMP WITH TIME ZONE): レコード作成日時。デフォルトは現在時刻。

### `users` テーブル (将来対応)
- **説明**: ユーザー情報を格納する。（現在は未実装）
- **カラム**:
    - `id` (UUID)
    - `email` (VARCHAR)
    - `name` (VARCHAR)
    - `created_at` (TIMESTAMP WITH TIME ZONE)

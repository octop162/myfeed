# カスタムフィードアプリ

## 概要

RSSフィードに加えて、カスタムスクレイピングによる多様なサイトからの情報収集が可能なフィードアプリケーションです。プログラマブルなプラグイン方式により、柔軟なデータ取得を実現します。

## システムアーキテクチャ

### 技術スタック

-   **バックエンド**: Go
-   **データベース**: PostgreSQL
-   **フロントエンド**: Next.js (React) - SPA構成
-   **API**: RESTful API
-   **コンテナ**: Docker

### アーキテクチャ図

```
[フロントエンド (Next.js)]
       ↕ RESTful API
[バックエンド (Go)]
       ↕
[データベース (PostgreSQL)]
       ↑
[バッチ処理 (スクレイピング)]
       ↑
[プラグインシステム]
```

## ローカル開発環境のセットアップ

プロジェクトをローカルで開発するために、以下の手順に従ってください。

### 前提条件

-   **Docker**: Docker Desktop または Docker Engine がインストールされていること。

### 手順

1.  **リポジトリをクローンする**:

    ```bash
    git clone <repository_url>
    cd feedapp
    ```

2.  **Docker Composeでアプリケーションを起動する**:

    プロジェクトのルートディレクトリで、以下のコマンドを実行します。

    ```bash
    docker-compose up --build
    ```

    これにより、PostgreSQLデータベース、Goバックエンド、Next.jsフロントエンドのすべてのサービスがビルドされ、起動します。

3.  **アプリケーションへのアクセス**:

    Webブラウザを開き、`http://localhost:3000` にアクセスします。

4.  **アプリケーションを停止する**:

    ```bash
    docker-compose down
    ```

## プロジェクト構造

```
/
├── cmd/                    # アプリケーションのエントリーポイント
│   └── server/             # バックエンドサーバーのエントリー
├── internal/               # プライベートアプリケーションコード
│   ├── config/             # 設定管理
│   ├── handler/            # HTTPハンドラー
│   ├── model/              # データモデル
│   ├── repository/         # データアクセス層
│   └── service/            # ビジネスロジック
├── migrations/             # データベースマイグレーションファイル
├── frontend/               # Next.jsフロントエンドアプリケーション
│   ├── public/             # 静的アセット
│   └── src/                # フロントエンドのソースコード
│       ├── app/            # Next.jsのApp Routerページ
│       ├── components/     # 再利用可能なReactコンポーネント
│       └── services/       # APIサービスクライアント
├── api/                    # API定義 (例: OpenAPI仕様)
├── docs/                   # プロジェクトドキュメント
├── pkg/                    # 公開ライブラリ (もしあれば)
└── plugins/                # プラグインファイル
```

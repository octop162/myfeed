# カスタムフィードアプリ要件定義書

## 概要

RSS フィードに加えて、カスタムスクレイピングによる多様なサイトからの情報収集が可能なフィードアプリケーション。プログラマブルなプラグイン方式により、柔軟なデータ取得を実現する。

## システム構成

### 技術スタック

- **バックエンド**: Go + Gin Framework
- **データベース**: PostgreSQL
- **フロントエンド**: Next.js (React) - SPA 構成
- **API**: RESTful API
- **コンテナ**: Docker
- **開発環境**: VS Code + devcontainer

### アーキテクチャ

```
[フロントエンド (Next.js)]
       ↕ RESTful API
[バックエンド (Go+Gin)]
       ↕
[データベース (PostgreSQL)]
       ↑
[バッチ処理 (スクレイピング)]
       ↑
[プラグインシステム]
```

## 機能要件

### 1. フィード管理

- **フィードソース登録**: URL、プラグイン種別、更新間隔の設定
- **フォルダ管理**: フラット構造での分類機能
  - デフォルトフォルダ「未分類」を提供
  - 各フィードは 1 つのフォルダにのみ所属
- **データ永続化**: 記事データの永続保存

### 2. プラグインシステム

- **実装方式**: Go plugin package
- **デフォルトプラグイン**: RSS/Atom フィード対応
- **カスタムプラグイン**: プログラマブルなスクレイピング処理
- **プラグイン設定**: 設定ファイルによる動的ロード

### 3. バッチ処理

- **更新間隔**: デフォルト 6 時間（サイト毎に設定可能）
- **並行処理**: Go routine による効率的なスクレイピング
- **エラーハンドリング**: ログ出力のみ（再試行なし）

### 4. 記事状態管理

- **4 つの状態パターン**:
  - 未読 + 通常
  - 未読 + 後で見る
  - 既読 + 通常
  - 既読 + 後で見る
- **後で見る**: 永続保存（自動削除なし）
- **重複検出**: 実装しない

### 5. Web UI

- **画面構成**:
  - 左サイドバー: フォルダ一覧
  - 右メインエリア: 記事一覧
- **表示形式**: シンプルなリスト形式
- **操作**: 記事の状態変更（未読/既読、後で見る切り替え）

### 6. API 設計

- **RESTful API**: 標準的な HTTP メソッド使用
- **エンドポイント例**:
  - `GET /api/folders` - フォルダ一覧
  - `GET /api/folders/{id}/articles` - フォルダ内記事一覧
  - `PUT /api/articles/{id}/status` - 記事状態更新
  - `GET /api/articles/later` - 後で見る記事一覧

## 非機能要件

### 1. 認証・認可

- **Phase 1**: 認証なし（ローカル環境想定）
- **Phase 2**: Google アカウント認証
- **マルチユーザー**: ユーザー毎のデータ分離

### 2. パフォーマンス

- **同時スクレイピング**: 制限なし
- **データベース**: PostgreSQL の全文検索機能は使用しない
- **フロントエンド**: SPA 構成によるレスポンシブな操作

### 3. 運用・監視

- **ログレベル**:
  - 開発環境: DEBUG
  - 本番環境: WARN, ERROR
- **監視**: 基本的なアプリケーションログのみ
- **通知機能**: 実装しない

### 4. セキュリティ

- **スクレイピング**: 対象サイトへの配慮は特に実装しない
- **データ保護**: PostgreSQL の標準的なセキュリティ機能を使用

## データモデル

### 主要エンティティ

```
Users (将来対応)
├── id, email, name, created_at

Folders
├── id, name, user_id, created_at

Feeds
├── id, name, url, plugin_type, folder_id,
├── update_interval, last_updated, created_at

Articles
├── id, feed_id, title, content, url,
├── published_at, is_read, is_later, created_at

Plugins
├── id, name, file_path, enabled, created_at
```

## 開発計画

### Phase 1: 基本機能

1. データベーススキーマ設計・構築
2. Go バックエンド API 開発
3. RSS プラグイン実装
4. Next.js フロントエンド基本画面
5. バッチ処理システム

### Phase 2: プラグインシステム

1. Go plugin package 実装
2. カスタムプラグインサンプル作成
3. プラグイン動的ロード機能

### Phase 3: 認証機能

1. Google OAuth 実装
2. マルチユーザー対応
3. データ分離機能

## 開発環境セットアップ

### devcontainer 構成

```dockerfile
# Go 1.21+
# Node.js 18+
# PostgreSQL 15+
# VS Code拡張機能
#   - Go
#   - PostgreSQL
#   - Docker
```

### 起動手順

1. `devcontainer`でコンテナ起動
2. PostgreSQL 初期化
3. Go モジュールダウンロード
4. Next.js 依存関係インストール
5. 開発サーバー起動

## コーディング規約

### Go (バックエンド)

#### コードスタイル

- **フォーマッター**: `gofmt`、`goimports`を使用
- **リンター**: `golangci-lint`を使用
- **命名規則**:
  - パッケージ名: 小文字、短縮形（例: `pkg`, `model`, `handler`）
  - 関数名: パスカルケース（例: `GetArticles`, `CreateFeed`）
  - 変数名: キャメルケース（例: `articleList`, `feedID`）
  - 定数: 大文字スネークケース（例: `DEFAULT_UPDATE_INTERVAL`）

#### ディレクトリ構成

```
/
├── cmd/                    # アプリケーションエントリーポイント
│   └── server/
├── internal/               # プライベートコード
│   ├── handler/           # HTTPハンドラー
│   ├── service/           # ビジネスロジック
│   ├── model/             # データモデル
│   ├── repository/        # データアクセス層
│   └── plugin/            # プラグインシステム
├── pkg/                   # 公開ライブラリ
├── api/                   # API定義（OpenAPI仕様書）
├── migrations/            # データベースマイグレーション
└── plugins/               # プラグインファイル
```

#### エラーハンドリング

- `errors.New()` または `fmt.Errorf()` を使用
- カスタムエラー型の定義
- ログレベルに応じた適切なログ出力

#### テスト

- **テスト駆動開発（TDD）**: Red-Green-Refactor サイクルを採用
- テストファイル: `*_test.go`
- テスト関数: `TestXxx(t *testing.T)`
- モック: `testify/mock` を使用
- カバレッジ目標: 80%以上
- **テスト種別**:
  - ユニットテスト: 各関数・メソッドの単体テスト
  - 統合テスト: API エンドポイントのテスト
  - E2E テスト: フロントエンドの主要フローテスト

### TypeScript/React (フロントエンド)

#### コードスタイル

- **フォーマッター**: Prettier
- **リンター**: ESLint + TypeScript ESLint
- **命名規則**:
  - コンポーネント: パスカルケース（例: `ArticleList`, `FeedItem`）
  - 関数・変数: キャメルケース（例: `fetchArticles`, `isLoading`）
  - ファイル名: kebab-case（例: `article-list.tsx`, `feed-item.tsx`）
  - 定数: 大文字スネークケース（例: `API_BASE_URL`）

#### ディレクトリ構成

```
src/
├── components/            # 再利用可能コンポーネント
│   ├── ui/               # UIコンポーネント
│   └── layout/           # レイアウトコンポーネント
├── pages/                # ページコンポーネント（Next.js）
├── hooks/                # カスタムフック
├── services/             # API呼び出し
├── types/                # TypeScript型定義
├── utils/                # ユーティリティ関数
└── styles/               # スタイルファイル
```

#### React 規則

- 関数コンポーネントを使用
- カスタムフックでロジック分離
- PropTypes ではなく TypeScript の型を使用
- デフォルトエクスポートよりも named export を推奨

#### 状態管理

- ローカル状態: `useState`, `useReducer`
- グローバル状態: Context API（必要に応じて Zustand）
- サーバー状態: SWR または React Query

### データベース

#### 命名規則

- **テーブル名**: 複数形、スネークケース（例: `users`, `feed_articles`）
- **カラム名**: スネークケース（例: `created_at`, `is_read`）
- **インデックス名**: `idx_テーブル名_カラム名`
- **外部キー名**: `fk_テーブル名_参照テーブル名`

#### マイグレーション

- ファイル名: `YYYYMMDDHHMMSS_description.sql`
- 前方互換性を保つ
- ロールバック用の DOWN スクリプトも作成

### API 設計

#### RESTful 原則

- **リソース指向**: 名詞で URL 設計
- **HTTP メソッド**: GET, POST, PUT, DELETE の適切な使用
- **ステータスコード**: HTTP 標準ステータスコードの遵守

#### URL 命名規則

```
GET    /api/v1/folders                 # フォルダ一覧取得
POST   /api/v1/folders                 # フォルダ作成
GET    /api/v1/folders/{id}            # フォルダ詳細取得
PUT    /api/v1/folders/{id}            # フォルダ更新
DELETE /api/v1/folders/{id}            # フォルダ削除
GET    /api/v1/folders/{id}/feeds      # フォルダ内フィード一覧
```

#### レスポンス形式

```json
{
  "data": {...},           # 成功時のデータ
  "message": "Success",    # メッセージ
  "timestamp": "2025-06-28T10:00:00Z"
}

{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力値が不正です",
    "details": [...]
  },
  "timestamp": "2025-06-28T10:00:00Z"
}
```

### プロジェクト管理・開発プロセス

#### テスト駆動開発（TDD）

プロジェクト全体でテスト駆動開発を採用し、以下のサイクルで開発を進める：

1. **Red**: 失敗するテストを先に書く
2. **Green**: テストをパスする最小限のコードを実装
3. **Refactor**: コードの品質を向上させる

#### GitHub Issue 管理

AI による開発支援を活用し、GitHub Issue でタスク管理を行う：

**Issue 作成規則**

- **タイトル**: `[Type] 簡潔な説明`
- **Type 種別**:
  - `Feature`: 新機能開発
  - `Bug`: バグ修正
  - `Refactor`: リファクタリング
  - `Test`: テスト追加・修正
  - `Docs`: ドキュメント更新

**Issue テンプレート例**

```markdown
## 概要

機能/修正の概要を記載

## 受け入れ条件

- [ ] 条件 1
- [ ] 条件 2
- [ ] テストが追加されている
- [ ] ドキュメントが更新されている

## 技術的詳細

実装時の注意点や参考資料

## 関連 Issue

- Related to #xxx
- Closes #xxx
```

**ラベル管理**

- `priority:high/medium/low`: 優先度
- `size:small/medium/large`: 作業量
- `component:api/frontend/database`: 対象コンポーネント
- `ai-assisted`: AI 支援による開発

**マイルストーン設定**

- Phase 1: 基本機能（MVP）
- Phase 2: プラグインシステム
- Phase 3: 認証機能

#### AI 開発支援ワークフロー

1. **Issue 作成**: GitHub Issue で要件定義
2. **AI 相談**: Claude 等の AI で技術的アプローチを検討
3. **TDD 実装**: テストファースト開発
4. **タスク消化と Issue 更新**: 実装したタスクは GitHub Issue を更新し、完了した項目にチェックを入れる。
5. **レビュー**: AI によるコードレビュー支援
6. **ドキュメント更新**: AI による文書作成支援

### Git 運用

#### ブランチ戦略

- **main**: 本番リリース用
- **develop**: 開発統合用
- **feature/xxx**: 機能開発用
- **hotfix/xxx**: 緊急修正用

#### コミットメッセージ

```
type(scope): subject

body

footer
```

**Type 種別**:

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: コードスタイル
- `refactor`: リファクタリング
- `test`: テスト
- `chore`: その他

#### コミット例

```
feat(api): フォルダ管理APIの実装

- フォルダCRUD操作の追加
- バリデーション機能の実装
- テストケースの追加

Closes #123
```

### 開発ツール設定

#### VS Code 拡張機能

```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-typescript-next",
    "ms-vscode.vscode-eslint",
    "ckolkman.vscode-postgres"
  ]
}
```

#### 自動フォーマット設定

```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

## 今後の拡張可能性

- **モバイルアプリ**: RESTful API 活用
- **エクスポート機能**: OPML, JSON 形式
- **テーマ機能**: ダーク/ライトモード
- **プラグインマーケット**: 共有可能なプラグインエコシステム
- **AI 要約機能**: 記事要約の自動生成

## AI 支援開発 TODO 計画

### 事前準備フェーズ

- [x] **開発環境セットアップ**

  - [x] devcontainer 設定ファイル作成（AI 支援で Dockerfile 作成）
  - [x] PostgreSQL 初期化スクリプト作成
  - [x] VS Code 設定ファイル（settings.json, extensions.json）作成
  - [x] Git 設定（.gitignore, commit template）

- [x] **プロジェクト初期化**
  - [x] Go mod ファイル初期化とディレクトリ構造作成
  - [x] Next.js プロジェクト初期化
  - [x] GitHub リポジトリ作成と Issue テンプレート設定
  - [x] CI/CD 設定（GitHub Actions）の基本構成

### Phase 1: 基本機能開発（AI 支援集中フェーズ）

#### データベース設計・実装

- [ ] **Issue #1: データベーススキーマ設計**
  - [x] AIとER図設計相談
  - [x] マイグレーションファイル作成（AIでSQL生成）
  - [x] テーブル関連のテストデータ作成

#### バックエンド API 開発

- [x] **Issue #2: プロジェクト構造とベースコード**

  - [x] Ginフレームワーク初期設定（AI支援でボイラープレート生成）
  - [x] ミドルウェア設定（CORS, ログ, エラーハンドリング）
  - [x] 設定管理システム実装

- [x] **Issue #3: フォルダ管理API（TDD実装）**

  - [x] テストケース作成（AI支援でエッジケース洗い出し）
  - [x] Handler実装
  - [x] Service層実装
  - [x] Repository層実装

- [x] **Issue #4: フィード管理API（TDD実装）**

  - [x] フィードCRUD操作のテスト作成
  - [x] API実装
  - [x] バリデーション機能実装

- [x] **Issue #5: 記事管理API（TDD実装）**
  - [x] 記事状態管理のテスト作成
  - [x] 記事取得・更新API実装
  - [ ] ページネーション実装

#### プラグインシステム基盤

- [ ] **Issue #6: プラグインアーキテクチャ設計**

  - [ ] AI とプラグインインターフェース設計相談
  - [ ] Go plugin パッケージ調査・実装方針決定
  - [ ] プラグイン設定ファイル形式決定

- [ ] **Issue #7: RSS プラグイン実装**
  - [ ] RSS パーサーのテスト作成（AI 支援でテストデータ生成）
  - [ ] RSS 取得・解析ロジック実装
  - [ ] エラーハンドリング実装

#### バッチ処理システム

- [ ] **Issue #8: スケジューラー実装**
  - [ ] cron 風スケジューラーの設計・実装
  - [ ] 並行処理での安全性確保
  - [ ] エラー処理とログ出力

#### フロントエンド基本実装

- [ ] **Issue #9: Next.js 基本セットアップ**

  - [ ] TypeScript 設定と ESLint/Prettier 設定
  - [ ] TailwindCSS 設定
  - [ ] 基本的なレイアウトコンポーネント作成

- [ ] **Issue #10: API 通信層実装**

  - [ ] API クライアントクラス設計・実装（AI 支援で型定義生成）
  - [ ] SWR または React Query 設定
  - [ ] エラーハンドリング実装

- [ ] **Issue #11: フォルダ管理 UI**

  - [ ] サイドバーコンポーネント作成
  - [ ] フォルダ一覧表示
  - [ ] フォルダ操作（作成・編集・削除）

- [ ] **Issue #12: 記事一覧 UI**
  - [ ] 記事リストコンポーネント作成
  - [ ] 記事状態変更機能（未読/既読、後で見る）
  - [ ] 無限スクロールまたはページネーション

### Phase 2: プラグインシステム拡張

#### カスタムプラグイン開発

- [ ] **Issue #13: プラグイン SDK 設計**

  - [ ] AI と協力してプラグイン開発者向けドキュメント作成
  - [ ] サンプルプラグインの実装
  - [ ] プラグインテストツール作成

- [ ] **Issue #14: 人気サイト用プラグイン実装**
  - [ ] Hacker News プラグイン（AI 支援でスクレイピングロジック）
  - [ ] Reddit プラグイン
  - [ ] Qiita プラグイン

#### プラグイン管理機能

- [ ] **Issue #15: プラグイン管理 UI**
  - [ ] プラグイン一覧・有効化/無効化
  - [ ] プラグイン設定画面
  - [ ] プラグインエラー表示

### Phase 3: 認証・マルチユーザー対応

#### 認証システム実装

- [ ] **Issue #16: Google OAuth 実装**

  - [ ] バックエンド認証ミドルウェア
  - [ ] フロントエンド認証フロー
  - [ ] セッション管理

- [ ] **Issue #17: マルチユーザー対応**
  - [ ] ユーザー毎のデータ分離
  - [ ] 権限管理システム
  - [ ] ユーザー設定画面

### 継続的改善フェーズ

#### 品質向上

- [ ] **Issue #18: パフォーマンス最適化**

  - [ ] データベースクエリ最適化（AI 支援で実行計画分析）
  - [ ] フロントエンドバンドル最適化
  - [ ] 画像最適化・遅延読み込み

- [ ] **Issue #19: 監視・ログ強化**
  - [ ] 構造化ログ実装
  - [ ] メトリクス収集
  - [ ] アラート設定

#### ドキュメント整備

- [ ] **Issue #20: API 仕様書作成**
  - [ ] OpenAPI 仕様書生成（AI 支援）
  - [ ] プラグイン開発ドキュメント
  - [ ] 運用マニュアル作成

### AI 活用の重点項目

#### コード生成・レビュー支援

- [ ] 各機能実装時のボイラープレート生成
- [ ] テストケース網羅性チェック
- [ ] セキュリティ脆弱性チェック
- [ ] パフォーマンス改善提案

#### 設計・アーキテクチャ相談

- [ ] 複雑な要件の整理・分解
- [ ] 技術選択の比較検討
- [ ] リファクタリング戦略立案

#### ドキュメント作成支援

- [ ] README.md 自動生成
- [ ] API ドキュメント生成
- [ ] コメント・型定義の自動補完

#### 開発進捗に伴うドキュメント更新管理

AI 支援により以下のドキュメントを開発と並行して継続的に更新：

**開発環境構築ドキュメント**

- [ ] **初期セットアップ手順書**
  - [ ] devcontainer 設定の詳細化（AI 支援で環境固有の問題解決）
  - [ ] 依存関係のバージョン管理とトラブルシューティング
  - [ ] 新メンバー向けオンボーディングガイド
- [ ] **ローカル開発手順**
  - [ ] ホットリロード設定とデバッグ手順
  - [ ] テスト実行方法とカバレッジ確認手順
  - [ ] プラグイン開発環境のセットアップ

**データベース構成ドキュメント**

- [ ] **スキーマ設計書**
  - [ ] 実装進捗に合わせた ERD の更新（AI 支援で図表生成）
  - [ ] テーブル定義の詳細化（制約、インデックス、トリガー）
  - [ ] データ型選択の根拠とパフォーマンス考慮事項
- [ ] **マイグレーション管理**
  - [ ] バージョン管理戦略とロールバック手順
  - [ ] 本番環境適用時の注意事項
  - [ ] パフォーマンス影響度の事前評価方法

**ディレクトリ構成ドキュメント**

- [ ] **プロジェクト構造ガイド**
  - [ ] 実装進捗に応じたディレクトリ構成の詳細化
  - [ ] 各ディレクトリの責務と依存関係
  - [ ] 新機能追加時のファイル配置ルール
- [ ] **コンポーネント設計書**
  - [ ] フロントエンド・バックエンドの層別設計
  - [ ] プラグインアーキテクチャの実装詳細
  - [ ] 設定ファイルとビルド成果物の管理

**ドキュメント更新の自動化**

- [ ] **CI/CD でのドキュメント生成**
  - [ ] コードコメントからの API 仕様書自動生成
  - [ ] データベースから現在のスキーマ図自動生成
  - [ ] 依存関係グラフの自動更新
- [ ] **リアルタイム更新**
  - [ ] プルリクエスト時のドキュメント差分チェック
  - [ ] ブランチマージ後の関連ドキュメント自動更新
  - [ ] 破壊的変更の影響範囲自動検出

**継続的ドキュメント改善プロセス**

```
開発実装 → AI支援レビュー → ドキュメント更新 → 検証 → フィードバック
     ↑                                                    ↓
     ←←←←←←←← 改善提案・修正依頼 ←←←←←←←←←←←←←←←←←←←
```

- [ ] **週次ドキュメント品質レビュー**
  - [ ] AI 支援での整合性チェック
  - [ ] 実装との乖離検出と修正
  - [ ] 新参者視点での理解しやすさ評価

### マイルストーン目標

**Milestone 1 (4 週間目標)**: MVP 完成

- ローカル環境でのフィード読み込み・表示
- 基本的なフォルダ・記事管理
- RSS フィード対応

**Milestone 2 (8 週間目標)**: プラグインシステム

- カスタムプラグイン機能
- 3 つ以上のサイト対応プラグイン
- プラグイン管理 UI

**Milestone 3 (12 週間目標)**: 本格運用対応

- Google 認証実装
- マルチユーザー対応
- 監視・ログシステム

### リスク管理

#### 技術的リスク

- [ ] Go plugin の動的ロードに関する制約調査
- [ ] 大量データ処理時のパフォーマンス検証
- [ ] スクレイピング対象サイトの変更対応

#### スケジュールリスク

- [ ] AI 支援効率の見積もり精度向上
- [ ] 複雑な機能の分割粒度調整
- [ ] テストカバレッジ維持とスピードのバランス

---

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/mattn/go-sqlite3" // SQLiteドライバー
	_ "github.com/lib/pq" // PostgreSQLドライバー

	"feedapp/internal/config"
	"feedapp/internal/handler"
	"feedapp/internal/repository"
	"feedapp/internal/service"
)

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	db, err := newDBConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリの初期化
	folderRepo := repository.NewFolderRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// サービスの初期化
	folderService := service.NewFolderService(folderRepo)
	feedService := service.NewFeedService(feedRepo)
	articleService := service.NewArticleService(articleRepo)

	// ハンドラの初期化
	folderHandler := handler.NewFolderHandler(folderService)
	feedHandler := handler.NewFeedHandler(feedService)
	articleHandler := handler.NewArticleHandler(articleService)

	// 環境変数からGIN_MODEを読み込み、Ginのモードを設定
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	r := gin.New()

	// ロガーミドルウェア
	r.Use(gin.Logger())

	// リカバリーミドルウェア（パニックからの回復）
	r.Use(gin.Recovery())

	// CORSミドルウェアの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // フロントエンドのURLを許可
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // 24時間
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// ルーティングの設定
	v1 := r.Group("/api/v1")
	{
		v1.GET("/folders", folderHandler.GetAllFolders)
		v1.GET("/folders/:id", folderHandler.GetFolderByID)
		v1.POST("/folders", folderHandler.CreateFolder)
		v1.PUT("/folders/:id", folderHandler.UpdateFolder)
		v1.DELETE("/folders/:id", folderHandler.DeleteFolder)

		v1.GET("/feeds", feedHandler.GetAllFeeds)
		v1.GET("/feeds/:id", feedHandler.GetFeedByID)
		v1.POST("/feeds", feedHandler.CreateFeed)
		v1.PUT("/feeds/:id", feedHandler.UpdateFeed)
		v1.DELETE("/feeds/:id", feedHandler.DeleteFeed)

		v1.GET("/articles", articleHandler.GetAllArticles)
		v1.GET("/articles/:id", articleHandler.GetArticleByID)
		v1.PUT("/articles/:id/status", articleHandler.UpdateArticleStatus)
		v1.GET("/articles/later", articleHandler.GetLaterArticles)
	}

	log.Printf("Server starting on :%s", cfg.Server.Port) // 設定からポートを読み込む
	if err := r.Run(":" + cfg.Server.Port); err != nil { // 設定からポートを読み込む
		log.Fatalf("Server failed to start: %v", err)
	}
}

// newDBConnection はデータベース接続を確立します。
func newDBConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	var db *sql.DB
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open database: %v. Retrying... (%d/%d)", err, i+1, maxRetries)
			time.Sleep(time.Second * 5) // 5秒待機
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("Failed to connect to database: %v. Retrying... (%d/%d)", err, i+1, maxRetries)
			db.Close()
			time.Sleep(time.Second * 5) // 5秒待機
			continue
		}
		log.Println("Successfully connected to database!")
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after multiple retries: %w", err)
	}

	// マイグレーションの実行
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations はデータベースマイグレーションを実行します。
func runMigrations(db *sql.DB) error {
	// create_tables.sql の実行
	createTablesSQL, err := os.ReadFile("migrations/20250628080159_create_tables.sql")
	if err != nil {
		return fmt.Errorf("failed to read create_tables.sql: %w", err)
	}
	_, err = db.Exec(string(createTablesSQL))
	if err != nil {
		return fmt.Errorf("failed to execute create_tables.sql: %w", err)
	}
	log.Println("create_tables.sql executed successfully.")

	// insert_test_data.sql の実行
	insertTestDataSQL, err := os.ReadFile("migrations/20250628081751_insert_test_data.sql")
	if err != nil {
		return fmt.Errorf("failed to read insert_test_data.sql: %w", err)
	}
	_, err = db.Exec(string(insertTestDataSQL))
	if err != nil {
		return fmt.Errorf("failed to execute insert_test_data.sql: %w", err)
	}
	log.Println("insert_test_data.sql executed successfully.")

	return nil
}

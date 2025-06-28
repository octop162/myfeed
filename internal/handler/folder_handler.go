// Package handler はHTTPリクエストを処理するハンドラーを提供します。
//
// このパッケージには、RESTful APIエンドポイントの実装が含まれており、
// HTTPリクエストの受信、バリデーション、サービス層の呼び出し、
// レスポンスの生成を担当します。
package handler

import (
	"errors"
	"net/http"

	"feedapp/internal/model"
	"feedapp/internal/service"

	"github.com/gin-gonic/gin"
)

// FolderHandler はフォルダ関連のHTTPリクエストを処理します。
//
// このハンドラーはフォルダのCRUD操作に対応するRESTful APIエンドポイントを提供し、
// フォルダサービス層と連携してビジネスロジックを実行します。
// 全てのエンドポイントはJSON形式でリクエストとレスポンスを扱います。
//
// サポートするエンドポイント:
//   - GET /folders - 全フォルダ一覧取得
//   - GET /folders/{id} - 特定フォルダ取得
//   - POST /folders - 新規フォルダ作成
//   - PUT /folders/{id} - フォルダ更新
//   - DELETE /folders/{id} - フォルダ削除
type FolderHandler struct {
	folderService service.FolderService // フォルダサービスへの依存
}

// NewFolderHandler は新しい FolderHandler インスタンスを作成します。
//
// 依存性注入パターンを使用してフォルダサービスを受け取り、
// 新しいハンドラーインスタンスを初期化します。
//
// Parameters:
//   - s: フォルダサービスのインターフェース実装
//
// Returns:
//   - *FolderHandler: 初期化されたハンドラーインスタンス
func NewFolderHandler(s service.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: s,
	}
}

// GetAllFolders はすべてのフォルダを取得します。
//
// GET /folders エンドポイントの実装です。
// データベースに保存されているすべてのフォルダを取得し、JSON配列として返します。
//
// HTTP Response:
//   - 200 OK: フォルダ配列のJSON
//   - 500 Internal Server Error: サーバー内部エラー
//
// Parameters:
//   - c: Ginのコンテキスト
func (h *FolderHandler) GetAllFolders(c *gin.Context) {
	folders, err := h.folderService.GetAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folders"})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// GetFolderByID は指定されたIDのフォルダを取得します。
//
// GET /folders/{id} エンドポイントの実装です。
// URLパラメータから取得したIDに基づいて、特定のフォルダ情報を取得します。
//
// HTTP Response:
//   - 200 OK: フォルダオブジェクトのJSON
//   - 404 Not Found: 指定されたIDのフォルダが存在しない
//   - 500 Internal Server Error: サーバー内部エラー
//
// Parameters:
//   - c: Ginのコンテキスト（URLパラメータ "id" を含む）
func (h *FolderHandler) GetFolderByID(c *gin.Context) {
	id := c.Param("id")
	folder, err := h.folderService.GetFolderByID(id)
	if err != nil {
		if errors.Is(err, service.ErrFolderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folder"})
		return
	}
	c.JSON(http.StatusOK, folder)
}

// CreateFolder は新しいフォルダを作成します。
//
// POST /folders エンドポイントの実装です。
// リクエストボディからフォルダ情報を取得し、新しいフォルダを作成します。
// IDと作成日時は自動的に設定されます。
//
// Request Body:
//   - name: フォルダ名（必須）
//   - user_id: ユーザーID（省略可能、将来対応）
//
// HTTP Response:
//   - 201 Created: 作成されたフォルダオブジェクトのJSON
//   - 400 Bad Request: リクエストボディの形式が不正
//   - 500 Internal Server Error: サーバー内部エラー
//
// Parameters:
//   - c: Ginのコンテキスト（JSONボディを含む）
func (h *FolderHandler) CreateFolder(c *gin.Context) {
	var folder model.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	createdFolder, err := h.folderService.CreateFolder(folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}
	c.JSON(http.StatusCreated, createdFolder)
}

// UpdateFolder は指定されたIDのフォルダを更新します。
//
// PUT /folders/{id} エンドポイントの実装です。
// URLパラメータのIDで指定されたフォルダの情報を、
// リクエストボディの内容で更新します。
//
// Request Body:
//   - name: フォルダ名（必須）
//   - user_id: ユーザーID（省略可能、将来対応）
//
// HTTP Response:
//   - 200 OK: 更新されたフォルダオブジェクトのJSON
//   - 400 Bad Request: リクエストボディの形式が不正
//   - 404 Not Found: 指定されたIDのフォルダが存在しない
//   - 500 Internal Server Error: サーバー内部エラー
//
// Parameters:
//   - c: Ginのコンテキスト（URLパラメータ "id" とJSONボディを含む）
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	id := c.Param("id")
	var folder model.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	folder.ID = id
	updatedFolder, err := h.folderService.UpdateFolder(id, folder)
	if err != nil {
		if errors.Is(err, service.ErrFolderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}
	c.JSON(http.StatusOK, updatedFolder)
}

// DeleteFolder は指定されたIDのフォルダを削除します。
//
// DELETE /folders/{id} エンドポイントの実装です。
// URLパラメータから取得したIDに基づいて、指定されたフォルダを削除します。
// 削除が成功した場合は、レスポンスボディなしで204 No Contentを返します。
//
// HTTP Response:
//   - 204 No Content: 削除成功（レスポンスボディなし）
//   - 404 Not Found: 指定されたIDのフォルダが存在しない
//   - 500 Internal Server Error: サーバー内部エラー
//
// Parameters:
//   - c: Ginのコンテキスト（URLパラメータ "id" を含む）
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	id := c.Param("id")
	err := h.folderService.DeleteFolder(id)
	if err != nil {
		if errors.Is(err, service.ErrFolderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}
	c.Status(http.StatusNoContent)
}

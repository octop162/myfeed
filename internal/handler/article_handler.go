package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"feedapp/internal/service"
)

// ArticleHandler は記事関連のHTTPリクエストを処理します。
type ArticleHandler struct {
	articleService service.ArticleService
}

// NewArticleHandler は新しい ArticleHandler インスタンスを作成します。
func NewArticleHandler(s service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: s,
	}
}

// GetAllArticles はすべての記事を取得します。
//
//	@Summary		記事一覧取得
//	@Description	データベースに保存されているすべての記事を取得します
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		model.Article	"記事一覧"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/articles [get]
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	articles, err := h.articleService.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// GetArticleByID は指定されたIDの記事を取得します。
//
//	@Summary		記事詳細取得
//	@Description	指定されたIDの記事を取得します
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"記事ID (UUID)"
//	@Success		200	{object}	model.Article	"記事詳細"
//	@Failure		404	{object}	map[string]string	"記事が見つかりません"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/articles/{id} [get]
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article"})
		return
	}
	c.JSON(http.StatusOK, article)
}

// ArticleStatusRequest は記事の状態更新リクエストを表します。
type ArticleStatusRequest struct {
	IsRead  bool `json:"is_read" example:"true"`
	IsLater bool `json:"is_later" example:"false"`
}

// UpdateArticleStatus は記事の状態を更新します。
//
//	@Summary		記事状態更新
//	@Description	指定されたIDの記事の読了状態や後で読む状態を更新します
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"記事ID (UUID)"
//	@Param			status	body		ArticleStatusRequest	true	"記事状態情報"
//	@Success		200	{object}	model.Article				"更新された記事"
//	@Failure		400	{object}	map[string]interface{}		"リクエストボディの形式が不正"
//	@Failure		404	{object}	map[string]string			"記事が見つかりません"
//	@Failure		500	{object}	map[string]string			"サーバー内部エラー"
//	@Router			/articles/{id}/status [put]
func (h *ArticleHandler) UpdateArticleStatus(c *gin.Context) {
	id := c.Param("id")
	var req ArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updatedArticle, err := h.articleService.UpdateArticleStatus(id, req.IsRead, req.IsLater)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article status"})
		return
	}
	c.JSON(http.StatusOK, updatedArticle)
}

// GetLaterArticles は「後で見る」に設定された記事を取得します。
//
//	@Summary		後で読む記事一覧取得
//	@Description	「後で読む」に設定された記事の一覧を取得します
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		model.Article	"後で読む記事一覧"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/articles/later [get]
func (h *ArticleHandler) GetLaterArticles(c *gin.Context) {
	articles, err := h.articleService.GetLaterArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get later articles"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

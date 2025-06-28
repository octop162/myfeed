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
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	articles, err := h.articleService.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// GetArticleByID は指定されたIDの記事を取得します。
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

// UpdateArticleStatus は記事の状態を更新します。
func (h *ArticleHandler) UpdateArticleStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		IsRead  bool `json:"is_read"`
		IsLater bool `json:"is_later"`
	}
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
func (h *ArticleHandler) GetLaterArticles(c *gin.Context) {
	articles, err := h.articleService.GetLaterArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get later articles"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

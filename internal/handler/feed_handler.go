package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"feedapp/internal/model"
	"feedapp/internal/service"
)

// FeedHandler はフィード関連のHTTPリクエストを処理します。
type FeedHandler struct {
	feedService service.FeedService
}

// NewFeedHandler は新しい FeedHandler インスタンスを作成します。
func NewFeedHandler(s service.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: s,
	}
}

// GetAllFeeds はすべてのフィードを取得します。
//
//	@Summary		フィード一覧取得
//	@Description	データベースに保存されているすべてのフィードを取得します
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		model.Feed	"フィード一覧"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/feeds [get]
func (h *FeedHandler) GetAllFeeds(c *gin.Context) {
	feeds, err := h.feedService.GetAllFeeds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds"})
		return
	}
	c.JSON(http.StatusOK, feeds)
}

// GetFeedByID は指定されたIDのフィードを取得します。
//
//	@Summary		フィード詳細取得
//	@Description	指定されたIDのフィードを取得します
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"フィードID (UUID)"
//	@Success		200	{object}	model.Feed	"フィード詳細"
//	@Failure		404	{object}	map[string]string	"フィードが見つかりません"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/feeds/{id} [get]
func (h *FeedHandler) GetFeedByID(c *gin.Context) {
	id := c.Param("id")
	feed, err := h.feedService.GetFeedByID(id)
	if err != nil {
		if errors.Is(err, service.ErrFeedNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}
	c.JSON(http.StatusOK, feed)
}

// CreateFeed は新しいフィードを作成します。
//
//	@Summary		フィード作成
//	@Description	新しいフィードを作成します
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			feed	body		model.Feed	true	"フィード情報"
//	@Success		201		{object}	model.Feed	"作成されたフィード"
//	@Failure		400		{object}	map[string]interface{}	"リクエストボディの形式が不正"
//	@Failure		500		{object}	map[string]string	"サーバー内部エラー"
//	@Router			/feeds [post]
func (h *FeedHandler) CreateFeed(c *gin.Context) {
	var feed model.Feed
	if err := c.ShouldBindJSON(&feed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	createdFeed, err := h.feedService.CreateFeed(feed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feed"})
		return
	}
	c.JSON(http.StatusCreated, createdFeed)
}

// UpdateFeed は指定されたIDのフィードを更新します。
//
//	@Summary		フィード更新
//	@Description	指定されたIDのフィードを更新します
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"フィードID (UUID)"
//	@Param			feed	body		model.Feed	true	"更新するフィード情報"
//	@Success		200		{object}	model.Feed	"更新されたフィード"
//	@Failure		400		{object}	map[string]interface{}	"リクエストボディの形式が不正"
//	@Failure		404		{object}	map[string]string	"フィードが見つかりません"
//	@Failure		500		{object}	map[string]string	"サーバー内部エラー"
//	@Router			/feeds/{id} [put]
func (h *FeedHandler) UpdateFeed(c *gin.Context) {
	id := c.Param("id")
	var feed model.Feed
	if err := c.ShouldBindJSON(&feed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	feed.ID = id
	updatedFeed, err := h.feedService.UpdateFeed(id, feed)
	if err != nil {
		if errors.Is(err, service.ErrFeedNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feed"})
		return
	}
	c.JSON(http.StatusOK, updatedFeed)
}

// DeleteFeed は指定されたIDのフィードを削除します。
//
//	@Summary		フィード削除
//	@Description	指定されたIDのフィードを削除します
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"フィードID (UUID)"
//	@Success		204	"削除成功"
//	@Failure		404	{object}	map[string]string	"フィードが見つかりません"
//	@Failure		500	{object}	map[string]string	"サーバー内部エラー"
//	@Router			/feeds/{id} [delete]
func (h *FeedHandler) DeleteFeed(c *gin.Context) {
	id := c.Param("id")
	err := h.feedService.DeleteFeed(id)
	if err != nil {
		if errors.Is(err, service.ErrFeedNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feed"})
		return
	}
		c.Status(http.StatusNoContent)
}

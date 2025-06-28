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
func (h *FeedHandler) GetAllFeeds(c *gin.Context) {
	feeds, err := h.feedService.GetAllFeeds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds"})
		return
	}
	c.JSON(http.StatusOK, feeds)
}

// GetFeedByID は指定されたIDのフィードを取得します。
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

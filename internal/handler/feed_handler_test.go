package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"feedapp/internal/model"
	"feedapp/internal/service"
)

// MockFeedService は service.FeedService のモック実装です。
type MockFeedService struct {
	mock.Mock
}

func (m *MockFeedService) GetAllFeeds() ([]model.Feed, error) {
	args := m.Called()
	return args.Get(0).([]model.Feed), args.Error(1)
}

func (m *MockFeedService) GetFeedByID(id string) (model.Feed, error) {
	args := m.Called(id)
	return args.Get(0).(model.Feed), args.Error(1)
}

func (m *MockFeedService) CreateFeed(feed model.Feed) (model.Feed, error) {
	args := m.Called(feed)
	return args.Get(0).(model.Feed), args.Error(1)
}

func (m *MockFeedService) UpdateFeed(id string, feed model.Feed) (model.Feed, error) {
	args := m.Called(id, feed)
	return args.Get(0).(model.Feed), args.Error(1)
}

func (m *MockFeedService) DeleteFeed(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestFeedHandler_GetAllFeeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFeedService)
	handler := NewFeedHandler(mockService)

	// 正常系: フィードが複数件ある場合
	t.Run("should return all feeds", func(t *testing.T) {
		expectedFeeds := []model.Feed{
			{ID: "1", Name: "Feed 1", URL: "http://example.com/feed1", PluginType: "rss"},
			{ID: "2", Name: "Feed 2", URL: "http://example.com/feed2", PluginType: "rss"},
		}
		mockService.On("GetAllFeeds").Return(expectedFeeds, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFeeds(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFeeds []model.Feed
		err := json.Unmarshal(w.Body.Bytes(), &actualFeeds)
		assert.NoError(t, err)
		assert.Equal(t, expectedFeeds, actualFeeds)
		mockService.AssertExpectations(t)
	})

	// 正常系: フィードが0件の場合
	t.Run("should return empty array if no feeds", func(t *testing.T) {
		mockService.On("GetAllFeeds").Return([]model.Feed{}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFeeds(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, "[]", w.Body.String())
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetAllFeeds").Return([]model.Feed{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFeeds(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get feeds")
		mockService.AssertExpectations(t)
	})
}

func TestFeedHandler_GetFeedByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFeedService)
	handler := NewFeedHandler(mockService)

	// 正常系: フィードが見つかる場合
	t.Run("should return feed by ID", func(t *testing.T) {
		expectedFeed := model.Feed{ID: "1", Name: "Feed 1", URL: "http://example.com/feed1", PluginType: "rss"}
		mockService.On("GetFeedByID", "1").Return(expectedFeed, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.GetFeedByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFeed model.Feed
		err := json.Unmarshal(w.Body.Bytes(), &actualFeed)
		assert.NoError(t, err)
		assert.Equal(t, expectedFeed, actualFeed)
		mockService.AssertExpectations(t)
	})

	// 異常系: フィードが見つからない場合
	t.Run("should return 404 if feed not found", func(t *testing.T) {
		mockService.On("GetFeedByID", "nonexistent").Return(model.Feed{}, service.ErrFeedNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		handler.GetFeedByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Feed not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetFeedByID", "errorID").Return(model.Feed{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		handler.GetFeedByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get feed")
		mockService.AssertExpectations(t)
	})
}

func TestFeedHandler_CreateFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFeedService)
	handler := NewFeedHandler(mockService)

	// 正常系: フィード作成成功
	t.Run("should create feed successfully", func(t *testing.T) {
		newFeed := model.Feed{Name: "New Feed", URL: "http://example.com/newfeed", PluginType: "rss", FolderID: "folder1"}
		createdFeed := model.Feed{ID: "3", Name: "New Feed", URL: "http://example.com/newfeed", PluginType: "rss", FolderID: "folder1"}
		mockService.On("CreateFeed", newFeed).Return(createdFeed, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"New Feed","url":"http://example.com/newfeed","plugin_type":"rss","folder_id":"folder1"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFeed(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var actualFeed model.Feed
		err := json.Unmarshal(w.Body.Bytes(), &actualFeed)
		assert.NoError(t, err)
		assert.Equal(t, createdFeed, actualFeed)
		mockService.AssertExpectations(t)
	})

	// 異常系: リクエストボディが不正（URLなし）
	t.Run("should return 400 if invalid request body (no url)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"New Feed","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFeed(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: リクエストボディが不正（不正なURL）
	t.Run("should return 400 if invalid request body (invalid url)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"New Feed","url":"invalid-url","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFeed(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		newFeed := model.Feed{Name: "Error Feed", URL: "http://example.com/errorfeed", PluginType: "rss"}
		mockService.On("CreateFeed", newFeed).Return(model.Feed{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"Error Feed","url":"http://example.com/errorfeed","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFeed(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create feed")
		mockService.AssertExpectations(t)
	})
}

func TestFeedHandler_UpdateFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFeedService)
	handler := NewFeedHandler(mockService)

	// 正常系: フィード更新成功
	t.Run("should update feed successfully", func(t *testing.T) {
		updatedFeed := model.Feed{ID: "1", Name: "Updated Feed", URL: "http://example.com/updatedfeed", PluginType: "rss", FolderID: "folder1"}
		mockService.On("UpdateFeed", "1", updatedFeed).Return(updatedFeed, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Updated Feed","url":"http://example.com/updatedfeed","plugin_type":"rss","folder_id":"folder1"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFeed(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFeed model.Feed
		err := json.Unmarshal(w.Body.Bytes(), &actualFeed)
		assert.NoError(t, err)
		assert.Equal(t, updatedFeed, actualFeed)
		mockService.AssertExpectations(t)
	})

	// 異常系: フィードが見つからない場合
	t.Run("should return 404 if feed not found", func(t *testing.T) {
		updatedFeed := model.Feed{ID: "nonexistent", Name: "Updated Feed", URL: "http://example.com/updatedfeed", PluginType: "rss"}
		mockService.On("UpdateFeed", "nonexistent", updatedFeed).Return(model.Feed{}, service.ErrFeedNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Updated Feed","url":"http://example.com/updatedfeed","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFeed(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Feed not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: リクエストボディが不正
	t.Run("should return 400 if invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Updated Feed","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFeed(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		updatedFeed := model.Feed{ID: "errorID", Name: "Error Feed", URL: "http://example.com/errorfeed", PluginType: "rss"}
		mockService.On("UpdateFeed", "errorID", updatedFeed).Return(model.Feed{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Error Feed","url":"http://example.com/errorfeed","plugin_type":"rss"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFeed(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update feed")
		mockService.AssertExpectations(t)
	})
}

func TestFeedHandler_DeleteFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFeedService)
	handler := NewFeedHandler(mockService)

	// 正常系: フィード削除成功
	t.Run("should delete feed successfully", func(t *testing.T) {
		mockService.On("DeleteFeed", "1").Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.DeleteFeed(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 異常系: フィードが見つからない場合
	t.Run("should return 404 if feed not found", func(t *testing.T) {
		mockService.On("DeleteFeed", "nonexistent").Return(service.ErrFeedNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		handler.DeleteFeed(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Feed not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("DeleteFeed", "errorID").Return(assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		handler.DeleteFeed(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete feed")
		mockService.AssertExpectations(t)
	})
}
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

// MockArticleService は service.ArticleService のモック実装です。
type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) GetAllArticles() ([]model.Article, error) {
	args := m.Called()
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *MockArticleService) GetArticleByID(id string) (model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) UpdateArticleStatus(id string, isRead, isLater bool) (model.Article, error) {
	args := m.Called(id, isRead, isLater)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) GetLaterArticles() ([]model.Article, error) {
	args := m.Called()
	return args.Get(0).([]model.Article), args.Error(1)
}

func TestArticleHandler_GetAllArticles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArticleService)
	handler := NewArticleHandler(mockService)

	// 正常系: 記事が複数件ある場合
	t.Run("should return all articles", func(t *testing.T) {
		expectedArticles := []model.Article{
			{ID: "1", Title: "Article 1", URL: "http://example.com/a1", IsRead: false, IsLater: false},
			{ID: "2", Title: "Article 2", URL: "http://example.com/a2", IsRead: true, IsLater: false},
		}
		mockService.On("GetAllArticles").Return(expectedArticles, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllArticles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualArticles []model.Article
		err := json.Unmarshal(w.Body.Bytes(), &actualArticles)
		assert.NoError(t, err)
		assert.Equal(t, expectedArticles, actualArticles)
		mockService.AssertExpectations(t)
	})

	// 正常系: 記事が0件の場合
	t.Run("should return empty array if no articles", func(t *testing.T) {
		mockService.On("GetAllArticles").Return([]model.Article{}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllArticles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, "[]", w.Body.String())
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetAllArticles").Return([]model.Article{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllArticles(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get articles")
		mockService.AssertExpectations(t)
	})
}

func TestArticleHandler_GetArticleByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArticleService)
	handler := NewArticleHandler(mockService)

	// 正常系: 記事が見つかる場合
	t.Run("should return article by ID", func(t *testing.T) {
		expectedArticle := model.Article{ID: "1", Title: "Article 1", URL: "http://example.com/a1", IsRead: false, IsLater: false}
		mockService.On("GetArticleByID", "1").Return(expectedArticle, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.GetArticleByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualArticle model.Article
		err := json.Unmarshal(w.Body.Bytes(), &actualArticle)
		assert.NoError(t, err)
		assert.Equal(t, expectedArticle, actualArticle)
		mockService.AssertExpectations(t)
	})

	// 異常系: 記事が見つからない場合
	t.Run("should return 404 if article not found", func(t *testing.T) {
		mockService.On("GetArticleByID", "nonexistent").Return(model.Article{}, service.ErrArticleNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		handler.GetArticleByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Article not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetArticleByID", "errorID").Return(model.Article{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		handler.GetArticleByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get article")
		mockService.AssertExpectations(t)
	})
}

func TestArticleHandler_UpdateArticleStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArticleService)
	handler := NewArticleHandler(mockService)

	// 正常系: 記事ステータス更新成功 (既読)
	t.Run("should update article status to read", func(t *testing.T) {
		updatedArticle := model.Article{ID: "1", Title: "Article 1", URL: "http://example.com/a1", IsRead: true, IsLater: false}
		mockService.On("UpdateArticleStatus", "1", true, false).Return(updatedArticle, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"is_read":true,"is_later":false}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateArticleStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualArticle model.Article
		err := json.Unmarshal(w.Body.Bytes(), &actualArticle)
		assert.NoError(t, err)
		assert.Equal(t, updatedArticle, actualArticle)
		mockService.AssertExpectations(t)
	})

	// 正常系: 記事ステータス更新成功 (後で見る)
	t.Run("should update article status to later", func(t *testing.T) {
		updatedArticle := model.Article{ID: "1", Title: "Article 1", URL: "http://example.com/a1", IsRead: false, IsLater: true}
		mockService.On("UpdateArticleStatus", "1", false, true).Return(updatedArticle, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"is_read":false,"is_later":true}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateArticleStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualArticle model.Article
		err := json.Unmarshal(w.Body.Bytes(), &actualArticle)
		assert.NoError(t, err)
		assert.Equal(t, updatedArticle, actualArticle)
		mockService.AssertExpectations(t)
	})

	// 異常系: 記事が見つからない場合
	t.Run("should return 404 if article not found", func(t *testing.T) {
		mockService.On("UpdateArticleStatus", "nonexistent", mock.Anything, mock.Anything).Return(model.Article{}, service.ErrArticleNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"is_read":true,"is_later":false}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateArticleStatus(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Article not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: リクエストボディが不正
	t.Run("should return 400 if invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"is_read":"invalid"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateArticleStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("UpdateArticleStatus", "errorID", mock.Anything, mock.Anything).Return(model.Article{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"is_read":true,"is_later":false}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateArticleStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update article status")
		mockService.AssertExpectations(t)
	})
}

func TestArticleHandler_GetLaterArticles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArticleService)
	handler := NewArticleHandler(mockService)

	// 正常系: 後で見る記事が複数件ある場合
	t.Run("should return later articles", func(t *testing.T) {
		expectedArticles := []model.Article{
			{ID: "1", Title: "Later Article 1", URL: "http://example.com/l1", IsRead: false, IsLater: true},
			{ID: "2", Title: "Later Article 2", URL: "http://example.com/l2", IsRead: true, IsLater: true},
		}
		mockService.On("GetLaterArticles").Return(expectedArticles, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetLaterArticles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualArticles []model.Article
		err := json.Unmarshal(w.Body.Bytes(), &actualArticles)
		assert.NoError(t, err)
		assert.Equal(t, expectedArticles, actualArticles)
		mockService.AssertExpectations(t)
	})

	// 正常系: 後で見る記事が0件の場合
	t.Run("should return empty array if no later articles", func(t *testing.T) {
		mockService.On("GetLaterArticles").Return([]model.Article{}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetLaterArticles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, "[]", w.Body.String())
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetLaterArticles").Return([]model.Article{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetLaterArticles(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get later articles")
		mockService.AssertExpectations(t)
	})
}
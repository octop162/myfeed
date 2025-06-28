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

// MockFolderService は service.FolderService のモック実装です。
type MockFolderService struct {
	mock.Mock
}

func (m *MockFolderService) GetAllFolders() ([]model.Folder, error) {
	args := m.Called()
	return args.Get(0).([]model.Folder), args.Error(1)
}

func (m *MockFolderService) GetFolderByID(id string) (model.Folder, error) {
	args := m.Called(id)
	return args.Get(0).(model.Folder), args.Error(1)
}

func (m *MockFolderService) CreateFolder(folder model.Folder) (model.Folder, error) {
	args := m.Called(folder)
	return args.Get(0).(model.Folder), args.Error(1)
}

func (m *MockFolderService) UpdateFolder(id string, folder model.Folder) (model.Folder, error) {
	args := m.Called(id, folder)
	return args.Get(0).(model.Folder), args.Error(1)
}

func (m *MockFolderService) DeleteFolder(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestFolderHandler_GetAllFolders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFolderService)
	handler := NewFolderHandler(mockService)

	// 正常系: フォルダが複数件ある場合
	t.Run("should return all folders", func(t *testing.T) {
		expectedFolders := []model.Folder{
			{ID: "1", Name: "Folder 1"},
			{ID: "2", Name: "Folder 2"},
		}
		mockService.On("GetAllFolders").Return(expectedFolders, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFolders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFolders []model.Folder
		err := json.Unmarshal(w.Body.Bytes(), &actualFolders)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolders, actualFolders)
		mockService.AssertExpectations(t)
	})

	// 正常系: フォルダが0件の場合
	t.Run("should return empty array if no folders", func(t *testing.T) {
		mockService.On("GetAllFolders").Return([]model.Folder{}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFolders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, "[]", w.Body.String())
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetAllFolders").Return([]model.Folder{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.GetAllFolders(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get folders")
		mockService.AssertExpectations(t)
	})
}

func TestFolderHandler_GetFolderByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFolderService)
	handler := NewFolderHandler(mockService)

	// 正常系: フォルダが見つかる場合
	t.Run("should return folder by ID", func(t *testing.T) {
		expectedFolder := model.Folder{ID: "1", Name: "Folder 1"}
		mockService.On("GetFolderByID", "1").Return(expectedFolder, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.GetFolderByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFolder model.Folder
		err := json.Unmarshal(w.Body.Bytes(), &actualFolder)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolder, actualFolder)
		mockService.AssertExpectations(t)
	})

	// 異常系: フォルダが見つからない場合
	t.Run("should return 404 if folder not found", func(t *testing.T) {
		mockService.On("GetFolderByID", "nonexistent").Return(model.Folder{}, service.ErrFolderNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		handler.GetFolderByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Folder not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("GetFolderByID", "errorID").Return(model.Folder{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		handler.GetFolderByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get folder")
		mockService.AssertExpectations(t)
	})
}

func TestFolderHandler_CreateFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFolderService)
	handler := NewFolderHandler(mockService)

	// 正常系: フォルダ作成成功
	t.Run("should create folder successfully", func(t *testing.T) {
		newFolder := model.Folder{Name: "New Folder"}
		createdFolder := model.Folder{ID: "3", Name: "New Folder"}
		mockService.On("CreateFolder", newFolder).Return(createdFolder, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"New Folder"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFolder(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var actualFolder model.Folder
		err := json.Unmarshal(w.Body.Bytes(), &actualFolder)
		assert.NoError(t, err)
		assert.Equal(t, createdFolder, actualFolder)
		mockService.AssertExpectations(t)
	})

	// 異常系: リクエストボディが不正
	t.Run("should return 400 if invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":""}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFolder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		newFolder := model.Folder{Name: "Error Folder"}
		mockService.On("CreateFolder", newFolder).Return(model.Folder{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"Error Folder"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateFolder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create folder")
		mockService.AssertExpectations(t)
	})
}

func TestFolderHandler_UpdateFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFolderService)
	handler := NewFolderHandler(mockService)

	// 正常系: フォルダ更新成功
	t.Run("should update folder successfully", func(t *testing.T) {
		updatedFolder := model.Folder{ID: "1", Name: "Updated Folder"}
		mockService.On("UpdateFolder", "1", updatedFolder).Return(updatedFolder, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Updated Folder"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFolder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualFolder model.Folder
		err := json.Unmarshal(w.Body.Bytes(), &actualFolder)
		assert.NoError(t, err)
		assert.Equal(t, updatedFolder, actualFolder)
		mockService.AssertExpectations(t)
	})

	// 異常系: フォルダが見つからない場合
	t.Run("should return 404 if folder not found", func(t *testing.T) {
		updatedFolder := model.Folder{ID: "nonexistent", Name: "Updated Folder"}
		mockService.On("UpdateFolder", "nonexistent", updatedFolder).Return(model.Folder{}, service.ErrFolderNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Updated Folder"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFolder(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Folder not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: リクエストボディが不正
	t.Run("should return 400 if invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":""}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFolder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		updatedFolder := model.Folder{ID: "errorID", Name: "Error Folder"}
		mockService.On("UpdateFolder", "errorID", updatedFolder).Return(model.Folder{}, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"name":"Error Folder"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateFolder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update folder")
		mockService.AssertExpectations(t)
	})
}

func TestFolderHandler_DeleteFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFolderService)
	handler := NewFolderHandler(mockService)

	// 正常系: フォルダ削除成功
	t.Run("should delete folder successfully", func(t *testing.T) {
		mockService.On("DeleteFolder", "1").Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		handler.DeleteFolder(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 異常系: フォルダが見つからない場合
	t.Run("should return 404 if folder not found", func(t *testing.T) {
		mockService.On("DeleteFolder", "nonexistent").Return(service.ErrFolderNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
		handler.DeleteFolder(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Folder not found")
		mockService.AssertExpectations(t)
	})

	// 異常系: サービスエラー
	t.Run("should return 500 if service error", func(t *testing.T) {
		mockService.On("DeleteFolder", "errorID").Return(assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "errorID"}}
		handler.DeleteFolder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete folder")
		mockService.AssertExpectations(t)
	})
}
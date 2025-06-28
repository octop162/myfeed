package handler

import (
	"errors"
	"net/http"

	"feedapp/internal/model"
	"feedapp/internal/service"

	"github.com/gin-gonic/gin"
)

// FolderHandler はフォルダ関連のHTTPリクエストを処理します。
type FolderHandler struct {
	folderService service.FolderService
}

// NewFolderHandler は新しい FolderHandler インスタンスを作成します。
func NewFolderHandler(s service.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: s,
	}
}

// GetAllFolders はすべてのフォルダを取得します。
func (h *FolderHandler) GetAllFolders(c *gin.Context) {
	folders, err := h.folderService.GetAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folders"})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// GetFolderByID は指定されたIDのフォルダを取得します。
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

package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"feedapp/internal/model"
)

var ErrNotFound = errors.New("not found")

// FolderRepository はフォルダのデータ永続化を定義するインターフェースです。
type FolderRepository interface {
	GetAll() ([]model.Folder, error)
	GetByID(id string) (model.Folder, error)
	Create(folder model.Folder) (model.Folder, error)
	Update(folder model.Folder) (model.Folder, error)
	Delete(id string) error
}

// folderRepository は FolderRepository インターフェースの実装です。
type folderRepository struct {
	db *sql.DB
}

// NewFolderRepository は新しい folderRepository インスタンスを作成します。
func NewFolderRepository(db *sql.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) GetAll() ([]model.Folder, error) {
	rows, err := r.db.Query("SELECT id, name, user_id, created_at FROM folders")
	if err != nil {
		return nil, fmt.Errorf("failed to get all folders: %w", err)
	}
	defer rows.Close()

	var folders []model.Folder
	for rows.Next() {
		var folder model.Folder
		var userID sql.NullString // user_idはNULL許容のためsql.NullStringを使用
		if err := rows.Scan(&folder.ID, &folder.Name, &userID, &folder.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan folder row: %w", err)
		}
		if userID.Valid {
			folder.UserID = userID.String
		}
		folders = append(folders, folder)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return folders, nil
}

func (r *folderRepository) GetByID(id string) (model.Folder, error) {
	var folder model.Folder
	var userID sql.NullString
	err := r.db.QueryRow("SELECT id, name, user_id, created_at FROM folders WHERE id = $1", id).Scan(&folder.ID, &folder.Name, &userID, &folder.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Folder{}, ErrNotFound
		}
		return model.Folder{}, fmt.Errorf("failed to get folder by ID: %w", err)
	}
	if userID.Valid {
		folder.UserID = userID.String
	}
	return folder, nil
}

func (r *folderRepository) Create(folder model.Folder) (model.Folder, error) {
		query := `INSERT INTO folders (id, name, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, name, user_id, created_at`
	var createdFolder model.Folder
	var userID sql.NullString
	err := r.db.QueryRow(query, folder.ID, folder.Name, sql.NullString{String: folder.UserID, Valid: folder.UserID != ""}, folder.CreatedAt).Scan(&createdFolder.ID, &createdFolder.Name, &userID, &createdFolder.CreatedAt)
	if err != nil {
		return model.Folder{}, fmt.Errorf("failed to create folder: %w", err)
	}
	if userID.Valid {
		createdFolder.UserID = userID.String
	}
	return createdFolder, nil
}

func (r *folderRepository) Update(folder model.Folder) (model.Folder, error) {
	query := `UPDATE folders SET name = $1 WHERE id = $2 RETURNING id, name, user_id, created_at`
	var updatedFolder model.Folder
	var userID string
	err := r.db.QueryRow(query, folder.Name, folder.ID).Scan(&updatedFolder.ID, &updatedFolder.Name, &userID, &updatedFolder.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Folder{}, ErrNotFound
		}
		return model.Folder{}, fmt.Errorf("failed to update folder: %w", err)
	}
	updatedFolder.UserID = userID
	return updatedFolder, nil
}

func (r *folderRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM folders WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}



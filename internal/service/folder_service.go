package service

import (
	"errors"
	"time" // timeパッケージを追加

	"feedapp/internal/model"
	"feedapp/internal/repository" // repositoryパッケージをインポート
)

var ErrFolderNotFound = errors.New("folder not found")

// FolderService はフォルダ関連のビジネスロジックを定義するインターフェースです。
type FolderService interface {
	GetAllFolders() ([]model.Folder, error)
	GetFolderByID(id string) (model.Folder, error)
	CreateFolder(folder model.Folder) (model.Folder, error)
	UpdateFolder(id string, folder model.Folder) (model.Folder, error)
	DeleteFolder(id string) error
}

// folderService は FolderService インターフェースの実装です。
type folderService struct {
	folderRepo repository.FolderRepository // Repositoryへの依存を追加
}

// NewFolderService は新しい folderService インスタンスを作成します。
func NewFolderService(repo repository.FolderRepository) FolderService {
	return &folderService{
		folderRepo: repo,
	}
}

func (s *folderService) GetAllFolders() ([]model.Folder, error) {
	folders, err := s.folderRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (s *folderService) GetFolderByID(id string) (model.Folder, error) {
	folder, err := s.folderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) { // repository.ErrNotFound を使用
			return model.Folder{}, ErrFolderNotFound
		}
		return model.Folder{}, err
	}
	return folder, nil
}

func (s *folderService) CreateFolder(folder model.Folder) (model.Folder, error) {
	// IDと作成日時をService層で設定
	folder.ID = model.GenerateUUID() // UUID生成関数を呼び出す
	folder.CreatedAt = time.Now()

	createdFolder, err := s.folderRepo.Create(folder)
	if err != nil {
		return model.Folder{}, err
	}
	return createdFolder, nil
}

func (s *folderService) UpdateFolder(id string, folder model.Folder) (model.Folder, error) {
	// 存在チェック
	_, err := s.folderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Folder{}, ErrFolderNotFound
		}
		return model.Folder{}, err
	}

	folder.ID = id // IDを設定
	updatedFolder, err := s.folderRepo.Update(folder)
	if err != nil {
		return model.Folder{}, err
	}
	return updatedFolder, nil
}

func (s *folderService) DeleteFolder(id string) error {
	err := s.folderRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFolderNotFound
		}
		return err
	}
	return nil
}

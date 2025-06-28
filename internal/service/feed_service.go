package service

import (
	"errors"
	"feedapp/internal/model"
	"feedapp/internal/repository"
	"time"
)

var ErrFeedNotFound = errors.New("feed not found")

// FeedService はフィード関連のビジネスロジックを定義するインターフェースです。
type FeedService interface {
	GetAllFeeds() ([]model.Feed, error)
	GetFeedByID(id string) (model.Feed, error)
	CreateFeed(feed model.Feed) (model.Feed, error)
	UpdateFeed(id string, feed model.Feed) (model.Feed, error)
	DeleteFeed(id string) error
}

// feedService は FeedService インターフェースの実装です。
type feedService struct {
	feedRepo repository.FeedRepository
}

// NewFeedService は新しい feedService インスタンスを作成します。
func NewFeedService(repo repository.FeedRepository) FeedService {
	return &feedService{
		feedRepo: repo,
	}
}

func (s *feedService) GetAllFeeds() ([]model.Feed, error) {
	feeds, err := s.feedRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (s *feedService) GetFeedByID(id string) (model.Feed, error) {
	feed, err := s.feedRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Feed{}, ErrFeedNotFound
		}
		return model.Feed{}, err
	}
	return feed, nil
}

func (s *feedService) CreateFeed(feed model.Feed) (model.Feed, error) {
	feed.ID = model.GenerateUUID()
	feed.CreatedAt = time.Now()
	createdFeed, err := s.feedRepo.Create(feed)
	if err != nil {
		return model.Feed{}, err
	}
	return createdFeed, nil
}

func (s *feedService) UpdateFeed(id string, feed model.Feed) (model.Feed, error) {
	_, err := s.feedRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Feed{}, ErrFeedNotFound
		}
		return model.Feed{}, err
	}
	feed.ID = id
	updatedFeed, err := s.feedRepo.Update(feed)
	if err != nil {
		return model.Feed{}, err
	}
	return updatedFeed, nil
}

func (s *feedService) DeleteFeed(id string) error {
	err := s.feedRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFeedNotFound
		}
		return err
	}
	return nil
}

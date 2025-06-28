package service

import (
	"errors"
	"feedapp/internal/model"
	"feedapp/internal/repository"
)

var ErrArticleNotFound = errors.New("article not found")

// ArticleService は記事関連のビジネスロジックを定義するインターフェースです。
type ArticleService interface {
	GetAllArticles() ([]model.Article, error)
	GetArticleByID(id string) (model.Article, error)
	UpdateArticleStatus(id string, isRead, isLater bool) (model.Article, error)
	GetLaterArticles() ([]model.Article, error)
}

// articleService は ArticleService インターフェースの実装です。
type articleService struct {
	articleRepo repository.ArticleRepository
}

// NewArticleService は新しい articleService インスタンスを作成します。
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepo: repo,
	}
}

func (s *articleService) GetAllArticles() ([]model.Article, error) {
	articles, err := s.articleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (s *articleService) GetArticleByID(id string) (model.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Article{}, ErrArticleNotFound
		}
		return model.Article{}, err
	}
	return article, nil
}

func (s *articleService) UpdateArticleStatus(id string, isRead, isLater bool) (model.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Article{}, ErrArticleNotFound
		}
		return model.Article{}, err
	}

	article.IsRead = isRead
	article.IsLater = isLater

	updatedArticle, err := s.articleRepo.Update(article)
	if err != nil {
		return model.Article{}, err
	}
	return updatedArticle, nil
}

func (s *articleService) GetLaterArticles() ([]model.Article, error) {
	articles, err := s.articleRepo.GetLaterArticles()
	if err != nil {
		return nil, err
	}
	return articles, nil
}

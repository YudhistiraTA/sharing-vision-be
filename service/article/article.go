package article

import (
	"context"

	"github.com/YudhistiraTA/sharing-vision-be/model"
)

type ArticleRepository interface {
	FindByID(ctx context.Context, id int) (model.Article, error)
	FindAll(ctx context.Context, limit int, offset int, status string) ([]model.Article, error)
	Create(ctx context.Context, article *model.Article) error
	Update(ctx context.Context, id int, article *model.Article) error
	Delete(ctx context.Context, id int) error
}

type ArticleService struct {
	Repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{Repo: repo}
}

func (s *ArticleService) FindByID(ctx context.Context, id int) (model.Article, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *ArticleService) FindAll(ctx context.Context, limit int, offset int, status string) ([]model.Article, error) {
	return s.Repo.FindAll(ctx, limit, offset, status)
}

func (s *ArticleService) Create(ctx context.Context, article *model.Article) error {
	return s.Repo.Create(ctx, article)
}

func (s *ArticleService) Update(ctx context.Context, id int, article *model.Article) error {
	return s.Repo.Update(ctx, id, article)
}

func (s *ArticleService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

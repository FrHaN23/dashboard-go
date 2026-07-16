package usecase

import (
	"context"

	"github.com/frhan23/dashboard-go/internal/entity"
	mysql "github.com/frhan23/dashboard-go/internal/module/article/repository"
)

type PaginationResult struct {
	Data       []entity.Article `json:"data"`
	TotalRows  int              `json:"total_rows"`
	TotalPages int              `json:"total_pages"`
}

type ArticleUsecase interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id int) (entity.Article, error)
	Fetch(ctx context.Context, limit, offset int, status string) (PaginationResult, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id int) error
}

type articleUsecase struct {
	repo mysql.ArticleRepository
}

func NewArticleUsecase(r mysql.ArticleRepository) ArticleUsecase {
	return &articleUsecase{
		repo: r,
	}
}

func (u *articleUsecase) Create(ctx context.Context, article *entity.Article) error {
	return u.repo.Create(ctx, article)
}

func (u *articleUsecase) GetByID(ctx context.Context, id int) (entity.Article, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *articleUsecase) Fetch(ctx context.Context, limit, offset int, status string) (PaginationResult, error) {
	articles, totalRows, err := u.repo.Fetch(ctx, limit, offset, status)
	if err != nil {
		return PaginationResult{}, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = totalRows / limit
		if totalRows%limit != 0 {
			totalPages++
		}
	}

	return PaginationResult{
		Data:       articles,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}, nil
}

func (u *articleUsecase) Update(ctx context.Context, article *entity.Article) error {
	return u.repo.Update(ctx, article)
}

func (u *articleUsecase) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

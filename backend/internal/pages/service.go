package pages

import (
	"context"
	"strings"

	"github.com/gosimple/slug"
	"github.com/iankencruz/threefive/internal/generated"
)

type PageService struct {
	Repo Repository
}

func NewPageService(repo Repository) *PageService {
	return &PageService{
		Repo: repo,
	}
}

// Create inserts a new page and auto-generates the slug from title
func (s *PageService) Create(ctx context.Context, arg generated.CreatePageParams) (*generated.Page, error) {
	arg.Slug = slug.Make(strings.ToLower(arg.Title))
	return s.Repo.CreatePage(ctx, arg)
}

// Update modifies a page and returns the fully updated page record
func (s *PageService) Update(ctx context.Context, arg generated.UpdatePageParams) (*generated.Page, error) {
	updatedPage, err := s.Repo.UpdatePage(ctx, arg)
	if err != nil {
		return nil, err
	}
	return updatedPage, nil
}

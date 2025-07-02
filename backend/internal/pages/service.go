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
	return &PageService{Repo: repo}
}

func (s *PageService) Create(ctx context.Context, arg generated.CreatePageParams) (*generated.Page, error) {
	arg.Slug = slug.Make(strings.ToLower(arg.Title))
	return s.Repo.CreatePage(ctx, arg)
}

func (s *PageService) Update(ctx context.Context, arg generated.UpdatePageParams) (*generated.Page, error) {
	return s.Repo.UpdatePage(ctx, arg)
}

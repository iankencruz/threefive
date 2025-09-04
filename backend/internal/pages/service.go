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

func (s *PageService) ListPages(ctx context.Context, sortParam string) ([]generated.Page, error) {
	field, direction := parseSortParam(sortParam)
	return s.Repo.ListPages(ctx, field, direction)
}

func parseSortParam(sort string) (field, direction string) {
	if sort == "" {
		return "updated_at", "desc"
	}

	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "updated_at", "desc"
	}

	field = strings.ToLower(parts[0])
	direction = strings.ToLower(parts[1])

	// validate allowed fields
	switch field {
	case "title", "created_at", "updated_at", "status":
		// ok
	default:
		field = "updated_at"
	}

	// validate direction
	if direction != "asc" && direction != "desc" {
		direction = "desc"
	}

	return field, direction
}

func (s *PageService) GetPublicPage(ctx context.Context, slug string) (PageWithGalleries, error) {
	return s.Repo.GetPublicPageWithGalleries(ctx, slug)
}

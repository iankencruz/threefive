package blogs

import (
	"context"

	"github.com/iankencruz/threefive/internal/generated"
)

type BlogService struct {
	repo Repository
}

func NewBlogService(repo Repository) *BlogService {
	return &BlogService{
		repo: repo,
	}
}

func (s *BlogService) GetGalleryBySlug(ctx context.Context, slug string) (BlogWithMedia, error) {
	blog, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return BlogWithMedia{}, err
	}

	mediaVals, err := s.repo.ListMediaForGallery(ctx, blog.ID)
	if err != nil {
		return BlogWithMedia{}, err
	}

	mediaPtrs := make([]*generated.Media, len(mediaVals))
	for i := range mediaVals {
		mediaPtrs[i] = &mediaVals[i]
	}

	return BlogWithMedia{
		Blog:  &blog,
		Media: mediaPtrs,
	}, nil
}

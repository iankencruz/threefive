package project

import (
	"context"
	"time"

	"github.com/iankencruz/threefive/internal/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectService struct {
	repo Repository
}

func NewProjectService(repo Repository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, arg generated.CreateProjectParams) (*generated.Project, error) {
	if arg.PublishedAt.Valid && !arg.IsPublished {
		arg.PublishedAt = pgtype.Timestamp{Valid: false}
	}
	return s.repo.CreateProject(ctx, arg)
}

func (s *ProjectService) Update(ctx context.Context, arg generated.UpdateProjectParams) (*generated.Project, error) {
	if arg.PublishedAt.Valid && !arg.IsPublished {
		arg.PublishedAt = pgtype.Timestamp{Valid: false}
	}

	var ts pgtype.Timestamp
	_ = ts.Scan(time.Now())
	arg.UpdatedAt = ts

	return s.repo.UpdateProject(ctx, arg)
}

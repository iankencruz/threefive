package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository interface {
	CreateProject(ctx context.Context, arg generated.CreateProjectParams) (*generated.Project, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*generated.Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*generated.Project, error)
	ListProjects(ctx context.Context) ([]generated.Project, error)
	UpdateProject(ctx context.Context, arg generated.UpdateProjectParams) (*generated.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error

	AddMediaToProject(ctx context.Context, arg generated.AddMediaToProjectParams) error
	RemoveMediaFromProject(ctx context.Context, arg generated.RemoveMediaFromProjectParams) error
	ListMediaForProject(ctx context.Context, projectID uuid.UUID) ([]generated.Media, error)
	UpdateProjectMediaSortOrder(ctx context.Context, arg generated.UpdateProjectMediaSortOrderParams) error
}

type ProjectRepository struct {
	q *generated.Queries
}

func NewRepository(q *generated.Queries) Repository {
	return &ProjectRepository{q: q}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, arg generated.CreateProjectParams) (*generated.Project, error) {
	project, err := r.q.CreateProject(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id uuid.UUID) (*generated.Project, error) {
	project, err := r.q.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetProjectBySlug(ctx context.Context, slug string) (*generated.Project, error) {
	project, err := r.q.GetProjectBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) ListProjects(ctx context.Context) ([]generated.Project, error) {
	return r.q.ListProjects(ctx)
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, arg generated.UpdateProjectParams) (*generated.Project, error) {
	project, err := r.q.UpdateProject(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProject(ctx, id)
}

func (r *ProjectRepository) AddMediaToProject(ctx context.Context, arg generated.AddMediaToProjectParams) error {
	return r.q.AddMediaToProject(ctx, arg)
}

func (r *ProjectRepository) RemoveMediaFromProject(ctx context.Context, arg generated.RemoveMediaFromProjectParams) error {
	return r.q.RemoveMediaFromProject(ctx, arg)
}

func (r *ProjectRepository) ListMediaForProject(ctx context.Context, projectID uuid.UUID) ([]generated.Media, error) {
	return r.q.ListMediaForProject(ctx, projectID)
}

func (r *ProjectRepository) UpdateProjectMediaSortOrder(ctx context.Context, arg generated.UpdateProjectMediaSortOrderParams) error {
	return r.q.UpdateProjectMediaSortOrder(ctx, arg)
}

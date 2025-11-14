package seo

import (
	"github.com/iankencruz/threefive/internal/shared/sqlc"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

// func (s *Service) UpsertSEO(ctx context.Context, qtx *sqlc.Queries, entityType string, entityID uuid.UUID, req *SEORequest) error {
// 	return errors.New("Upsert SEO")
// }

// func (s *Service) GetSEO(ctx context.Context, entityType string, entityID uuid.UUID) (*SEOResponse, error) {
// 	return errors.New("GetSEO")
// }
// func (s *Service) DeleteSEO(ctx context.Context, entityType string, entityID uuid.UUID) error

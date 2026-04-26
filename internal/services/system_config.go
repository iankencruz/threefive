package services

import (
	"context"
	"fmt"

	"github.com/iankencruz/threefive/database/generated"
)

//
// type SystemConfig struct {
// 	configCode string
// 	value      string
// 	created_at time.Time
// 	updated_at time.Time
// }

type SystemConfigService struct {
	queries *generated.Queries
}

func NewSystemConfigService(queries *generated.Queries) *SystemConfigService {
	return &SystemConfigService{
		queries: queries,
	}
}

// ListProjects retrieves a paginated list of projects
func (s *SystemConfigService) ListSystemConfig(ctx context.Context, limit, offset int32) ([]generated.SystemConfig, error) {
	configs, err := s.queries.ListConfig(ctx, generated.ListConfigParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return configs, nil
}

// ListProjects retrieves a paginated list of projects
func (s *SystemConfigService) GetConfigByCode(ctx context.Context, code string) (generated.SystemConfig, error) {
	config, err := s.queries.GetConfigByCode(ctx, code)
	if err != nil {
		return generated.SystemConfig{}, fmt.Errorf("failed to list configs: %w", err)
	}

	return config, nil
}

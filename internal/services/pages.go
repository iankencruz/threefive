// internal/services/page.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type SocialLinks struct {
	Twitter   string `json:"twitter,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
	GitHub    string `json:"github,omitempty"`
	Instagram string `json:"instagram,omitempty"`
}

// PageResponse - only used when we need to enrich page data with related entities
type PageResponse struct {
	Page             generated.Page
	HeroMedia        *MediaResponse
	ContentImage     *MediaResponse
	SocialLinks      *SocialLinks
	FeaturedProjects []FeaturedProjectSummary
}

type FeaturedProjectSummary struct {
	Project       generated.GetFeaturedProjectsRow
	FeaturedImage *MediaResponse
	GalleryMedia  []MediaResponse
}

type PageService struct {
	queries      *generated.Queries
	mediaService *MediaService
}

func NewPageService(queries *generated.Queries, mediaService *MediaService) *PageService {
	return &PageService{
		queries:      queries,
		mediaService: mediaService,
	}
}

// ListPages returns all pages
func (s *PageService) ListPages(ctx context.Context) ([]generated.Page, error) {
	pages, err := s.queries.ListPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}

	pageOrder := map[string]int{
		"home":    0,
		"about":   1,
		"contact": 2,
	}

	sort.Slice(pages, func(i, j int) bool {
		orderI, existsI := pageOrder[pages[i].Slug]
		orderJ, existsJ := pageOrder[pages[j].Slug]
		if existsI && existsJ {
			return orderI < orderJ
		}
		if existsI {
			return true
		}
		if existsJ {
			return false
		}
		return pages[i].Slug < pages[j].Slug
	})

	return pages, nil
}

// GetPageBySlug returns a page with enriched data (media loaded)
func (s *PageService) GetPageBySlug(ctx context.Context, slug string) (*PageResponse, error) {
	page, err := s.queries.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	return s.enrichPageData(ctx, &page)
}

// UpdatePageBySlug updates a page using form data
func (s *PageService) UpdatePageBySlug(ctx context.Context, slug string, params generated.UpdatePageParams) (*PageResponse, error) {
	// Get existing page to get its ID
	existing, err := s.queries.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Set ID in params
	params.ID = existing.ID

	// Update page
	updated, err := s.queries.UpdatePage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	// Return enriched response so the handler always has fresh data
	return s.enrichPageData(ctx, &updated)
}

// UpdateFeaturedProjects updates the featured projects for a page
func (s *PageService) UpdateFeaturedProjects(ctx context.Context, pageID pgtype.UUID, projectIDs []string) error {
	if len(projectIDs) > 3 {
		return fmt.Errorf("maximum 3 featured projects allowed")
	}

	// Clear existing
	if err := s.queries.ClearFeaturedProjects(ctx, pageID); err != nil {
		return fmt.Errorf("failed to clear featured projects: %w", err)
	}

	// Add new ones in order
	for i, projectIDStr := range projectIDs {
		projectUUID, err := uuid.Parse(projectIDStr)
		if err != nil {
			continue
		}

		_, err = s.queries.AddFeaturedProject(ctx, generated.AddFeaturedProjectParams{
			PageID:       pageID,
			ProjectID:    pgtype.UUID{Bytes: projectUUID, Valid: true},
			DisplayOrder: int32(i),
		})
		if err != nil {
			return fmt.Errorf("failed to add featured project: %w", err)
		}
	}

	return nil
}

// enrichPageData loads related media and projects for display
func (s *PageService) enrichPageData(ctx context.Context, page *generated.Page) (*PageResponse, error) {
	response := &PageResponse{
		Page: *page,
	}

	// Load hero media if set
	if page.HeroMediaID.Valid {
		heroMedia, err := s.mediaService.GetMediaByID(ctx, page.HeroMediaID)
		if err == nil {
			mediaResp := s.mediaService.ToMediaResponse(heroMedia)
			response.HeroMedia = &mediaResp
		}
	}

	// Load content image if set (for about page)
	if page.ContentImageID.Valid {
		contentImg, err := s.mediaService.GetMediaByID(ctx, page.ContentImageID)
		if err == nil {
			mediaResp := s.mediaService.ToMediaResponse(contentImg)
			response.ContentImage = &mediaResp
		}
	}

	// Parse social links JSON (for contact page)
	if len(page.SocialLinks) > 0 {
		var social SocialLinks
		if err := json.Unmarshal(page.SocialLinks, &social); err == nil {
			response.SocialLinks = &social
		}
	}

	// Load featured projects for home AND about pages
	if page.PageType == "home" || page.PageType == "about" {
		featuredProjects, err := s.loadFeaturedProjects(ctx, page.ID)
		if err == nil {
			response.FeaturedProjects = featuredProjects
		}
	}

	return response, nil
}

func (s *PageService) loadFeaturedProjects(ctx context.Context, pageID pgtype.UUID) ([]FeaturedProjectSummary, error) {
	rows, err := s.queries.GetFeaturedProjects(ctx, pageID)
	if err != nil {
		return nil, err
	}

	var projects []FeaturedProjectSummary
	for _, row := range rows {
		project := FeaturedProjectSummary{
			Project: row,
		}

		if row.FeaturedImageID.Valid {
			featuredImg, err := s.mediaService.GetMediaByID(ctx, row.FeaturedImageID)
			if err == nil {
				mediaResp := s.mediaService.ToMediaResponse(featuredImg)
				project.FeaturedImage = &mediaResp
			}
		}

		galleryMedia, err := s.mediaService.GetGalleryMediaForEntity(ctx, "project", row.ID)
		if err == nil && len(galleryMedia) > 0 {
			project.GalleryMedia = s.mediaService.ToMediaResponses(galleryMedia)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

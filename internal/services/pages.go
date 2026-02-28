// internal/services/pages.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

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
// Similar to how MediaResponse adds URLs to generated.Media
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

// UpdatePageRequest represents the data needed to update a page
type UpdatePageRequest struct {
	// Common fields for all pages
	Title       string
	Header      string
	SubHeader   string
	HeroMediaID string

	// About page specific
	Content            string
	ContentImageID     string
	CtaText            string
	CtaLink            string
	FeaturedProjectIDs []string

	// Contact page specific
	Email           string
	SocialTwitter   string
	SocialLinkedIn  string
	SocialGitHub    string
	SocialInstagram string
}

// Validate validates the update request based on page type
func (r *UpdatePageRequest) Validate(pageType string) (map[string]string, error) {
	errors := make(map[string]string)

	// Title is only required for contact and about pages (not home)
	if pageType != "home" && strings.TrimSpace(r.Title) == "" {
		errors["title"] = "Title is required"
	}

	// Header validation (optional, but if provided must not be empty)
	if strings.TrimSpace(r.Header) != "" && len(r.Header) > 500 {
		errors["header"] = "Header must be less than 500 characters"
	}

	// Sub-header validation (optional)
	if strings.TrimSpace(r.SubHeader) != "" && len(r.SubHeader) > 1000 {
		errors["sub_header"] = "Sub-header must be less than 1000 characters"
	}

	// Page-specific validation
	switch pageType {
	case "contact":
		if strings.TrimSpace(r.Email) == "" {
			errors["email"] = "Email is required for contact page"
		} else if !isValidEmail(r.Email) {
			errors["email"] = "Invalid email format"
		}

	case "about":
		// Content is optional but has max length
		if len(r.Content) > 5000 {
			errors["content"] = "Content must be less than 5000 characters"
		}
	}

	if len(errors) > 0 {
		return errors, fmt.Errorf("validation failed")
	}

	return nil, nil
}

// Helper function for email validation
func isValidEmail(email string) bool {
	// Simple email regex
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
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

	// Sort pages in specific order: home, about, contact, then any others
	pageOrder := map[string]int{
		"home":    0,
		"about":   1,
		"contact": 2,
	}

	sort.Slice(pages, func(i, j int) bool {
		orderI, existsI := pageOrder[pages[i].Slug]
		orderJ, existsJ := pageOrder[pages[j].Slug]

		// If both pages are in the order map, sort by their defined order
		if existsI && existsJ {
			return orderI < orderJ
		}
		// If only one is in the map, it comes first
		if existsI {
			return true
		}
		if existsJ {
			return false
		}
		// If neither is in the map, sort alphabetically
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

// UpdatePageBySlug updates a page using the UpdatePageRequest
func (s *PageService) UpdatePageBySlug(ctx context.Context, slug string, req *UpdatePageRequest) (*PageResponse, error) {
	// Get existing page to get its ID
	existing, err := s.queries.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Build SQLC params
	params := generated.UpdatePageParams{
		ID: existing.ID,
	}

	// Set fields from request
	if req.Title != "" {
		params.Title = pgtype.Text{String: req.Title, Valid: true}
	}
	if req.Header != "" {
		params.Header = pgtype.Text{String: req.Header, Valid: true}
	}
	if req.SubHeader != "" {
		params.SubHeader = pgtype.Text{String: req.SubHeader, Valid: true}
	}
	if req.HeroMediaID != "" {
		heroUUID, err := uuid.Parse(req.HeroMediaID)
		if err == nil {
			params.HeroMediaID = pgtype.UUID{Bytes: heroUUID, Valid: true}
		}
	}
	if req.Content != "" {
		params.Content = pgtype.Text{String: req.Content, Valid: true}
	}
	if req.ContentImageID != "" {
		contentUUID, err := uuid.Parse(req.ContentImageID)
		if err == nil {
			params.ContentImageID = pgtype.UUID{Bytes: contentUUID, Valid: true}
		}
	}
	if req.CtaText != "" {
		params.CtaText = pgtype.Text{String: req.CtaText, Valid: true}
	}
	if req.CtaLink != "" {
		params.CtaLink = pgtype.Text{String: req.CtaLink, Valid: true}
	}
	if req.Email != "" {
		params.Email = pgtype.Text{String: req.Email, Valid: true}
	}

	// Update page
	updated, err := s.queries.UpdatePage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	// ✅ UPDATE FEATURED PROJECTS (for both home and about pages)
	if len(req.FeaturedProjectIDs) >= 0 { // Allow clearing by passing empty array
		maxProjects := 6 // Home page allows 6
		if existing.PageType == "about" {
			maxProjects = 3 // About page allows 3
		}

		if len(req.FeaturedProjectIDs) > maxProjects {
			return nil, fmt.Errorf("maximum %d featured projects allowed", maxProjects)
		}

		if err := s.UpdateFeaturedProjects(ctx, existing.ID, req.FeaturedProjectIDs); err != nil {
			return nil, fmt.Errorf("failed to update featured projects: %w", err)
		}
	}

	// Return enriched page data
	return s.enrichPageData(ctx, &updated)
}

// UpdateFeaturedProjects updates the featured projects for a page (About page only)
func (s *PageService) UpdateFeaturedProjects(ctx context.Context, pageID pgtype.UUID, projectIDs []string) error {
	// Validate max 3 projects
	if len(projectIDs) > 3 {
		return fmt.Errorf("maximum 3 featured projects allowed")
	}

	// Clear existing
	if err := s.queries.ClearFeaturedProjects(ctx, pageID); err != nil {
		return fmt.Errorf("failed to clear featured projects: %w", err)
	}

	// Add new ones
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

	// Load featured projects for about page
	if page.PageType == "about" {
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

		// Load featured image if set
		if row.FeaturedImageID.Valid {
			featuredImg, err := s.mediaService.GetMediaByID(ctx, row.FeaturedImageID)
			if err == nil {
				mediaResp := s.mediaService.ToMediaResponse(featuredImg)
				project.FeaturedImage = &mediaResp
			}
		}

		// Load gallery media for this project
		galleryMedia, err := s.mediaService.GetGalleryMediaForEntity(ctx, "project", row.ID)
		if err == nil && len(galleryMedia) > 0 {
			project.GalleryMedia = s.mediaService.ToMediaResponses(galleryMedia)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

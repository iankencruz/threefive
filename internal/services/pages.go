// internal/services/pages.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/validation"
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
func (r *UpdatePageRequest) Validate(pageType string) (validation.FieldErrors, error) {
	fields := []validation.Field{
		{
			Name:  "title",
			Value: r.Title,
			Rules: []validation.ValidationRule{
				validation.Required("Title is required"),
				validation.MaxLength(200, "Title must be at most 200 characters"),
			},
		},
		{
			Name:  "header",
			Value: r.Header,
			Rules: []validation.ValidationRule{
				validation.MaxLength(500, "Header must be at most 500 characters"),
			},
		},
		{
			Name:  "sub_header",
			Value: r.SubHeader,
			Rules: []validation.ValidationRule{
				validation.MaxLength(500, "Sub-header must be at most 500 characters"),
			},
		},
	}

	// Page-specific validation
	switch pageType {
	case "about":
		fields = append(fields,
			validation.Field{
				Name:  "content",
				Value: r.Content,
				Rules: []validation.ValidationRule{
					validation.MaxLength(10000, "Content must be at most 10,000 characters"),
				},
			},
			validation.Field{
				Name:  "cta_text",
				Value: r.CtaText,
				Rules: []validation.ValidationRule{
					validation.MaxLength(100, "CTA text must be at most 100 characters"),
				},
			},
			validation.Field{
				Name:  "cta_link",
				Value: r.CtaLink,
				Rules: []validation.ValidationRule{
					validation.IsURL(""),
				},
			},
		)

	case "contact":
		fields = append(fields,
			validation.Field{
				Name:  "email",
				Value: r.Email,
				Rules: []validation.ValidationRule{
					validation.Required("Email is required for contact page"),
					validation.IsEmail(""),
				},
			},
			validation.Field{
				Name:  "social_twitter",
				Value: r.SocialTwitter,
				Rules: []validation.ValidationRule{
					validation.IsURL(""),
				},
			},
			validation.Field{
				Name:  "social_linkedin",
				Value: r.SocialLinkedIn,
				Rules: []validation.ValidationRule{
					validation.IsURL(""),
				},
			},
			validation.Field{
				Name:  "social_github",
				Value: r.SocialGitHub,
				Rules: []validation.ValidationRule{
					validation.IsURL(""),
				},
			},
			validation.Field{
				Name:  "social_instagram",
				Value: r.SocialInstagram,
				Rules: []validation.ValidationRule{
					validation.IsURL(""),
				},
			},
		)
	}

	errors := validation.ValidateFields(fields)
	if errors.HasErrors() {
		return errors, fmt.Errorf("validation failed")
	}

	return nil, nil
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

// UpdatePageBySlug updates a page using the new request struct
func (s *PageService) UpdatePageBySlug(ctx context.Context, slug string, req *UpdatePageRequest) (*PageResponse, error) {
	// Get existing page
	existing, err := s.queries.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Build update params
	params := generated.UpdatePageParams{
		ID: existing.ID,
	}

	// Common fields (all pages)
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
		if mediaUUID, err := uuid.Parse(req.HeroMediaID); err == nil {
			params.HeroMediaID = pgtype.UUID{Bytes: mediaUUID, Valid: true}
		}
	}

	// About page specific
	if existing.PageType == "about" {
		if req.Content != "" {
			params.Content = pgtype.Text{String: req.Content, Valid: true}
		}
		if req.ContentImageID != "" {
			if imgUUID, err := uuid.Parse(req.ContentImageID); err == nil {
				params.ContentImageID = pgtype.UUID{Bytes: imgUUID, Valid: true}
			}
		}
		if req.CtaText != "" {
			params.CtaText = pgtype.Text{String: req.CtaText, Valid: true}
		}
		if req.CtaLink != "" {
			params.CtaLink = pgtype.Text{String: req.CtaLink, Valid: true}
		}

		// Update featured projects
		if len(req.FeaturedProjectIDs) > 0 {
			if err := s.UpdateFeaturedProjects(ctx, existing.ID, req.FeaturedProjectIDs); err != nil {
				return nil, fmt.Errorf("failed to update featured projects: %w", err)
			}
		}
	}

	// Contact page specific
	if existing.PageType == "contact" {
		if req.Email != "" {
			params.Email = pgtype.Text{String: req.Email, Valid: true}
		}

		// Handle social links JSON
		socialLinks := SocialLinks{
			Twitter:   req.SocialTwitter,
			LinkedIn:  req.SocialLinkedIn,
			GitHub:    req.SocialGitHub,
			Instagram: req.SocialInstagram,
		}

		// Only set if at least one link exists
		if socialLinks.Twitter != "" || socialLinks.LinkedIn != "" || socialLinks.GitHub != "" || socialLinks.Instagram != "" {
			socialLinksJSON, _ := json.Marshal(socialLinks)
			params.SocialLinks = socialLinksJSON
		}
	}

	// Update page
	updated, err := s.queries.UpdatePage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	// Return enriched response
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

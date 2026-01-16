// internal/contacts/service.go
package contacts

import (
	"context"
	"log"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/email"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	emailService *email.Service
}

func NewService(db *pgxpool.Pool, queries *sqlc.Queries, emailService *email.Service) *Service {
	return &Service{
		db:           db,
		queries:      queries,
		emailService: emailService,
	}
}

// CreateContact creates a new contact submission
func (s *Service) CreateContact(ctx context.Context, req *ContactRequest, ipAddress, userAgent string) (*ContactResponse, error) {
	// Check rate limit (max 3 submissions per IP per hour)
	if err := s.checkRateLimit(ctx, ipAddress); err != nil {
		return nil, err
	}

	// Parse IP address
	var ipAddr *netip.Addr
	if ipAddress != "" {
		parsed, err := netip.ParseAddr(ipAddress)
		if err == nil {
			ipAddr = &parsed
		}
	}

	// Handle optional subject
	var subject pgtype.Text
	if req.Subject != "" {
		subject = pgtype.Text{String: req.Subject, Valid: true}
	}

	// Handle optional user agent
	var ua pgtype.Text
	if userAgent != "" {
		ua = pgtype.Text{String: userAgent, Valid: true}
	}

	contact, err := s.queries.CreateContact(ctx, sqlc.CreateContactParams{
		Name:      req.Name,
		Email:     req.Email,
		Subject:   subject,
		Message:   req.Message,
		IpAddress: ipAddr,
		UserAgent: ua,
	})
	if err != nil {
		return nil, errors.Internal("Failed to create contact", err)
	}

	// Send email notification (non-blocking, graceful degradation)
	go s.sendEmailNotification(contact)

	return s.toContactResponse(contact), nil
}

// sendEmailNotification sends email notification asynchronously
func (s *Service) sendEmailNotification(contact sqlc.Contacts) {
	ctx := context.Background()

	subjectText := "No subject"
	if contact.Subject.Valid {
		subjectText = contact.Subject.String
	}

	err := s.emailService.SendContactNotification(
		contact.Name,
		contact.Email,
		subjectText,
		contact.Message,
	)

	if err != nil {
		log.Printf("Failed to send email for contact %s: %v", contact.ID, err)
		// Mark email as failed
		s.queries.MarkEmailFailed(ctx, sqlc.MarkEmailFailedParams{
			ID:    contact.ID,
			Error: pgtype.Text{String: err.Error(), Valid: true},
		})
	} else {
		log.Printf("Email sent successfully for contact %s", contact.ID)
		// Mark email as sent
		s.queries.MarkEmailSent(ctx, contact.ID)
	}
}

// RetryFailedEmails attempts to resend emails that failed
func (s *Service) RetryFailedEmails(ctx context.Context) error {
	contacts, err := s.queries.GetUnsentEmails(ctx, 50) // Retry up to 50 at a time
	if err != nil {
		return err
	}

	log.Printf("Retrying %d failed emails", len(contacts))

	for _, contact := range contacts {
		subjectText := "No subject"
		if contact.Subject.Valid {
			subjectText = contact.Subject.String
		}

		err := s.emailService.SendContactNotification(
			contact.Name,
			contact.Email,
			subjectText,
			contact.Message,
		)

		if err != nil {
			log.Printf("Retry failed for contact %s: %v", contact.ID, err)
			s.queries.MarkEmailFailed(ctx, sqlc.MarkEmailFailedParams{
				ID:    contact.ID,
				Error: pgtype.Text{String: err.Error(), Valid: true},
			})
		} else {
			log.Printf("Retry successful for contact %s", contact.ID)
			s.queries.MarkEmailSent(ctx, contact.ID)
		}
	}

	return nil
}

// GetContactByID retrieves a contact by ID
func (s *Service) GetContactByID(ctx context.Context, id uuid.UUID) (*ContactResponse, error) {
	contact, err := s.queries.GetContactByID(ctx, id)
	if err != nil {
		return nil, errors.Internal("Failed to get contact", err)
	}

	return s.toContactResponse(contact), nil
}

// ListContacts retrieves contacts with pagination and filtering
func (s *Service) ListContacts(ctx context.Context, params ListContactsParams) (*ContactListResponse, error) {
	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.OrderBy == "" {
		params.OrderBy = "created_at_desc"
	}

	var contacts []sqlc.Contacts
	var total int64
	var err error

	// Get contacts based on filter
	if params.Status != "" {
		contacts, err = s.queries.ListContactsByStatus(ctx, sqlc.ListContactsByStatusParams{
			Status:      params.Status, // status is string, not pgtype.Text
			LimitCount:  params.Limit,
			OffsetCount: params.Offset,
		})
		if err != nil {
			return nil, errors.Internal("Failed to list contacts by status", err)
		}

		total, err = s.queries.CountContactsByStatus(ctx, params.Status) // string, not pgtype.Text
	} else {
		contacts, err = s.queries.ListContacts(ctx, sqlc.ListContactsParams{
			OrderBy:     params.OrderBy, // string, not pgtype.Text
			LimitCount:  params.Limit,
			OffsetCount: params.Offset,
		})
		if err != nil {
			return nil, errors.Internal("Failed to list contacts", err)
		}

		total, err = s.queries.CountContacts(ctx)
	}

	if err != nil {
		return nil, errors.Internal("Failed to count contacts", err)
	}

	// Convert to response
	contactResponses := make([]ContactResponse, len(contacts))
	for i, contact := range contacts {
		contactResponses[i] = *s.toContactResponse(contact)
	}

	totalPages := total / int64(params.Limit)
	if total%int64(params.Limit) > 0 {
		totalPages++
	}

	return &ContactListResponse{
		Contacts:   contactResponses,
		Total:      total,
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalPages: totalPages,
	}, nil
}

// UpdateContactStatus updates the status of a contact
func (s *Service) UpdateContactStatus(ctx context.Context, id uuid.UUID, status string) (*ContactResponse, error) {
	contact, err := s.queries.UpdateContactStatus(ctx, sqlc.UpdateContactStatusParams{
		Status: status, // status is string, not pgtype.Text
		ID:     id,
	})
	if err != nil {
		return nil, errors.Internal("Failed to update contact status", err)
	}

	return s.toContactResponse(contact), nil
}

// DeleteContact soft deletes a contact
func (s *Service) DeleteContact(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.SoftDeleteContact(ctx, id); err != nil {
		return errors.Internal("Failed to delete contact", err)
	}
	return nil
}

// checkRateLimit checks if the IP has exceeded the submission rate limit
func (s *Service) checkRateLimit(ctx context.Context, ipAddress string) error {
	if ipAddress == "" {
		return nil // Skip rate limit if no IP
	}

	ipAddr, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil // Skip if IP can't be parsed
	}

	// Check submissions in the last hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	contacts, err := s.queries.GetContactsByIPAddress(ctx, sqlc.GetContactsByIPAddressParams{
		IpAddress: &ipAddr,
		AfterTime: oneHourAgo,
	})
	if err != nil {
		// Don't fail on rate limit check error
		return nil
	}

	// Max 3 submissions per hour
	if len(contacts) >= 3 {
		return errors.BadRequest("Rate limit exceeded: maximum 3 submissions per hour", "rate_limit_exceeded")
	}

	return nil
}

// toContactResponse converts sqlc.Contacts to ContactResponse
func (s *Service) toContactResponse(contact sqlc.Contacts) *ContactResponse {
	resp := &ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Message:   contact.Message,
		Status:    contact.Status, // status is already string
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	if contact.Subject.Valid {
		resp.Subject = &contact.Subject.String
	}

	if contact.DeletedAt.Valid {
		resp.DeletedAt = &contact.DeletedAt.Time
	}

	return resp
}

// CleanupDeletedContacts removes contacts deleted more than 30 days ago
func (s *Service) CleanupDeletedContacts(ctx context.Context) error {
	if err := s.queries.CleanupDeletedContacts(ctx); err != nil {
		return errors.Internal("Failed to cleanup deleted contacts", err)
	}
	return nil
}

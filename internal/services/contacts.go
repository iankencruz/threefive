package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/validation"
	"github.com/jackc/pgx/v5/pgtype"
)

const maxEmailAttempts = 5

// ── Request ───────────────────────────────────────────────────────────────────

type ContactFormRequest struct {
	FirstName string
	LastName  string
	Email     string
	Subject   string
	Message   string
}

func (r *ContactFormRequest) Validate() (validation.FieldErrors, error) {
	fields := []validation.Field{
		{
			Name:  "first_name",
			Value: r.FirstName,
			Rules: []validation.ValidationRule{
				validation.Required("First name is required"),
				validation.MaxLength(100, ""),
			},
		},
		{
			Name:  "last_name",
			Value: r.LastName,
			Rules: []validation.ValidationRule{
				validation.Required("Last name is required"),
				validation.MaxLength(100, ""),
			},
		},
		{
			Name:  "email",
			Value: r.Email,
			Rules: []validation.ValidationRule{
				validation.Required("Email is required"),
				validation.IsEmail(""),
			},
		},
		{
			Name:  "subject",
			Value: r.Subject,
			Rules: []validation.ValidationRule{
				validation.Required("Subject is required"),
				validation.MaxLength(200, ""),
			},
		},
		{
			Name:  "message",
			Value: r.Message,
			Rules: []validation.ValidationRule{
				validation.Required("Message is required"),
				validation.MinLength(10, "Message must be at least 10 characters"),
			},
		},
	}

	errs := validation.ValidateFields(fields)
	if errs.HasErrors() {
		return errs, fmt.Errorf("validation failed")
	}
	return nil, nil
}

// ── Service ───────────────────────────────────────────────────────────────────

type ContactService struct {
	queries   *generated.Queries
	logger    *slog.Logger
	smtpHost  string
	smtpPort  string
	smtpUser  string
	smtpPass  string
	fromEmail string
	toEmail   string
}

func NewContactService(
	queries *generated.Queries,
	logger *slog.Logger,
	smtpHost, smtpPort, smtpUser, smtpPass, fromEmail, toEmail string,
) *ContactService {
	return &ContactService{
		queries:   queries,
		logger:    logger,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		smtpUser:  smtpUser,
		smtpPass:  smtpPass,
		fromEmail: fromEmail,
		toEmail:   toEmail,
	}
}

// Submit saves the submission then attempts to send the email immediately.
// The user always gets a success response as long as the DB save worked.
func (s *ContactService) Submit(ctx context.Context, req *ContactFormRequest) (*generated.ContactSubmission, error) {
	// 1. Save to DB
	submission, err := s.queries.CreateContactSubmission(ctx, generated.CreateContactSubmissionParams{
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Email:     strings.TrimSpace(req.Email),
		Subject:   strings.TrimSpace(req.Subject),
		Message:   strings.TrimSpace(req.Message),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save contact submission: %w", err)
	}

	// 2. Attempt email immediately
	if sendErr := s.sendEmail(req); sendErr != nil {
		s.logger.Warn("initial email send failed, will retry later",
			"submission_id", submission.ID,
			"error", sendErr,
		)
		// Record the failed attempt
		_ = s.queries.MarkEmailFailed(ctx, generated.MarkEmailFailedParams{
			ID:         submission.ID,
			EmailError: pgtype.Text{String: sendErr.Error(), Valid: true},
		})
	} else {
		// Mark as sent immediately
		_ = s.queries.MarkEmailSent(ctx, submission.ID)
		s.logger.Info("contact email sent successfully", "submission_id", submission.ID)
	}

	return &submission, nil
}

// RetryUnsent is called by the background worker to retry failed emails.
func (s *ContactService) RetryUnsent(ctx context.Context) {
	unsent, err := s.queries.GetUnsentSubmissions(ctx, maxEmailAttempts)
	if err != nil {
		s.logger.Error("failed to fetch unsent contact submissions", "error", err)
		return
	}

	if len(unsent) == 0 {
		return
	}

	s.logger.Info("retrying unsent contact emails", "count", len(unsent))

	for _, sub := range unsent {
		req := &ContactFormRequest{
			FirstName: sub.FirstName,
			LastName:  sub.LastName,
			Email:     sub.Email,
			Subject:   sub.Subject,
			Message:   sub.Message,
		}

		if sendErr := s.sendEmail(req); sendErr != nil {
			s.logger.Warn("retry email send failed",
				"submission_id", sub.ID,
				"attempts", sub.EmailAttempts+1,
				"error", sendErr,
			)
			_ = s.queries.MarkEmailFailed(ctx, generated.MarkEmailFailedParams{
				ID:         sub.ID,
				EmailError: pgtype.Text{String: sendErr.Error(), Valid: true},
			})
		} else {
			_ = s.queries.MarkEmailSent(ctx, sub.ID)
			s.logger.Info("retry email sent successfully",
				"submission_id", sub.ID,
				"attempts", sub.EmailAttempts+1,
			)
		}
	}
}

func (s *ContactService) sendEmail(req *ContactFormRequest) error {
	if s.smtpHost == "" || s.toEmail == "" {
		return nil // SMTP not configured, skip
	}

	subject := fmt.Sprintf("[Contact] %s", req.Subject)
	body := fmt.Sprintf(
		"Name: %s %s\r\nEmail: %s\r\n\r\n%s",
		req.FirstName, req.LastName, req.Email, req.Message,
	)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.fromEmail, s.toEmail, subject, body,
	)

	var auth smtp.Auth
	if s.smtpUser != "" && s.smtpPass != "" {
		auth = smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	}

	return smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.fromEmail, []string{s.toEmail}, []byte(msg))
}

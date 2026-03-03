package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/validation"
)

// ── Request ──────────────────────────────────────────────────────────────────

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
	smtpHost  string
	smtpPort  string
	smtpUser  string
	smtpPass  string
	fromEmail string
	toEmail   string
}

func NewContactService(
	queries *generated.Queries,
	smtpHost, smtpPort, smtpUser, smtpPass, fromEmail, toEmail string,
) *ContactService {
	return &ContactService{
		queries:   queries,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		smtpUser:  smtpUser,
		smtpPass:  smtpPass,
		fromEmail: fromEmail,
		toEmail:   toEmail,
	}
}

func (s *ContactService) Submit(ctx context.Context, req *ContactFormRequest) (*generated.ContactSubmission, error) {
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

	go func() {
		_ = s.sendEmail(req)
	}()

	return &submission, nil
}

func (s *ContactService) sendEmail(req *ContactFormRequest) error {
	if s.smtpHost == "" || s.toEmail == "" {
		return nil
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

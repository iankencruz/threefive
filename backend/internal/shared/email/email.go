// backend/internal/shared/email/email.go
package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/iankencruz/threefive/internal/config"
)

type Service struct {
	config config.SMTPConfig
	auth   smtp.Auth
}

func NewService(cfg config.SMTPConfig) *Service {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &Service{
		config: cfg,
		auth:   auth,
	}
}

// SendContactNotification sends an email notification for a new contact submission
func (s *Service) SendContactNotification(contactName, contactEmail, subject, message string) error {
	// Email template
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9fafb; padding: 20px; margin-top: 20px; border-radius: 8px; }
        .field { margin-bottom: 15px; }
        .label { font-weight: bold; color: #6B7280; }
        .value { margin-top: 5px; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #E5E7EB; color: #6B7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Contact Form Submission</h1>
        </div>
        
        <div class="content">
            <div class="field">
                <div class="label">From:</div>
                <div class="value">%s (%s)</div>
            </div>
            
            <div class="field">
                <div class="label">Subject:</div>
                <div class="value">%s</div>
            </div>
            
            <div class="field">
                <div class="label">Message:</div>
                <div class="value">%s</div>
            </div>
        </div>
        
        <div class="footer">
            This is an automated notification from your website contact form.
        </div>
    </div>
</body>
</html>
`, contactName, contactEmail, subject, template.HTMLEscapeString(message))

	// Plain text alternative
	textBody := fmt.Sprintf(`
New Contact Form Submission

From: %s (%s)
Subject: %s

Message:
%s

---
This is an automated notification from your website contact form.
`, contactName, contactEmail, subject, message)

	return s.SendEmail(s.config.From, "New Contact Form Submission", htmlBody, textBody)
}

// SendEmail sends an email with both HTML and plain text versions
func (s *Service) SendEmail(to, subject, htmlBody, textBody string) error {
	// Build email headers and body
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n")
	buf.WriteString("\r\n")

	// Plain text version
	buf.WriteString("--boundary123\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(textBody)
	buf.WriteString("\r\n")

	// HTML version
	buf.WriteString("--boundary123\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)
	buf.WriteString("\r\n")

	buf.WriteString("--boundary123--\r\n")

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, s.auth, s.config.From, []string{to}, buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

package handler

import (
	"log/slog"

	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/labstack/echo/v5"
)

type ContactHandler struct {
	logger         *slog.Logger
	contactService *services.ContactService
}

func NewContactHandler(logger *slog.Logger, contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{
		logger:         logger,
		contactService: contactService,
	}
}

func (h *ContactHandler) HandleSubmit(c *echo.Context) error {
	ctx := c.Request().Context()
	req := &services.ContactFormRequest{
		FirstName: c.FormValue("first_name"),
		LastName:  c.FormValue("last_name"),
		Email:     c.FormValue("email"),
		Subject:   c.FormValue("subject"),
		Message:   c.FormValue("message"),
	}

	// Validate
	fieldErrors, err := req.Validate()
	if err != nil {
		component := pages.ContactForm(req, fieldErrors)
		return responses.RenderError(c.Request().Context(), c, component, "Please fix the errors below")
	}

	// Submit
	if _, err := h.contactService.Submit(ctx, req); err != nil {
		h.logger.Error("failed to submit contact form", "error", err)
		return responses.Render(ctx, c, pages.ContactForm(req, map[string]string{
			"general": "Something went wrong. Please try again.",
		}))
	}

	// Tell HTMX to update the browser URL to /contact
	c.Response().Header().Set("HX-Push-Url", "/contact")

	return responses.Render(ctx, c, pages.ContactSuccess())
}

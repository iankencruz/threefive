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
	req := &services.ContactFormRequest{
		FirstName: c.FormValue("first_name"),
		LastName:  c.FormValue("last_name"),
		Email:     c.FormValue("email"),
		Subject:   c.FormValue("subject"),
		Message:   c.FormValue("message"),
	}

	fieldErrors, err := req.Validate()
	if err != nil {
		component := pages.ContactForm(req, fieldErrors)
		return responses.RenderError(c.Request().Context(), c, component, "Please fix the errors below")
	}

	_, err = h.contactService.Submit(c.Request().Context(), req)
	if err != nil {
		h.logger.Error("failed to submit contact form", "error", err)
		component := pages.ContactForm(req, map[string]string{
			"general": "Something went wrong. Please try again.",
		})
		return responses.RenderError(c.Request().Context(), c, component, "Submission failed")
	}

	return responses.Render(c.Request().Context(), c, pages.ContactSuccess())
}

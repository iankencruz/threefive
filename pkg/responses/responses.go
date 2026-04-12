// pkg/responses/responses.go
package responses

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/iankencruz/threefive/templates/components/toast"
	"github.com/labstack/echo/v5"
)

// Render renders a templ component
func Render(ctx context.Context, c *echo.Context, component templ.Component) error {
	return component.Render(ctx, c.Response())
}

// RenderWithStatus renders a component with a specific status code
func RenderWithStatus(ctx context.Context, c *echo.Context, status int, component templ.Component) error {
	c.Response().WriteHeader(status)
	return component.Render(ctx, c.Response())
}

// RenderWithToast renders a component with a toast notification using OOB swap
func RenderWithToast(ctx context.Context, c *echo.Context, component templ.Component, message string, variant toast.Variant) error {
	// Render main component first
	if err := component.Render(ctx, c.Response()); err != nil {
		return err
	}

	// Create toast component
	toastComponent := toast.Toast(toast.Props{
		Variant:       variant,
		Description:   message,
		Icon:          true,
		ShowIndicator: true,
	})

	// Wrap toast in OOB div for HTMX to swap into #toast-container
	oobToast := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Write OOB wrapper opening tag
		if _, err := w.Write([]byte(`<div hx-swap-oob="afterbegin:#toast-container">`)); err != nil {
			return err
		}

		// Render toast component
		if err := toastComponent.Render(ctx, w); err != nil {
			return err
		}

		// Write OOB wrapper closing tag
		if _, err := w.Write([]byte(`</div>`)); err != nil {
			return err
		}

		return nil
	})

	// Render the OOB toast
	return oobToast.Render(ctx, c.Response())
}

// RenderSuccess renders with success toast
func RenderSuccess(ctx context.Context, c *echo.Context, component templ.Component, message string) error {
	return RenderWithToast(ctx, c, component, message, toast.VariantSuccess)
}

// RenderError renders with error toast
func RenderError(ctx context.Context, c *echo.Context, component templ.Component, message string) error {
	return RenderWithToast(ctx, c, component, message, toast.VariantError)
}

// RenderWarning renders with warning toast
func RenderWarning(ctx context.Context, c *echo.Context, component templ.Component, message string) error {
	return RenderWithToast(ctx, c, component, message, toast.VariantDefault)
}

// RenderInfo renders with info toast
func RenderInfo(ctx context.Context, c *echo.Context, component templ.Component, message string) error {
	return RenderWithToast(ctx, c, component, message, toast.VariantDefault)
}

// ToastOnly renders only a toast notification with OOB swap
func ToastOnly(ctx context.Context, c *echo.Context, message string, variant toast.Variant) error {
	toastComponent := toast.Toast(toast.Props{
		Variant:       variant,
		Description:   message,
		Icon:          true,
		ShowIndicator: true,
	})

	// Wrap in OOB div
	oobToast := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := w.Write([]byte(`<div hx-swap-oob="afterbegin:#toast-container">`)); err != nil {
			return err
		}

		if err := toastComponent.Render(ctx, w); err != nil {
			return err
		}

		if _, err := w.Write([]byte(`</div>`)); err != nil {
			return err
		}

		return nil
	})

	return oobToast.Render(ctx, c.Response())
}

// SuccessToast renders only a success toast
func SuccessToast(ctx context.Context, c *echo.Context, message string) error {
	return ToastOnly(ctx, c, message, toast.VariantSuccess)
}

// ErrorToast renders only an error toast
func ErrorToast(ctx context.Context, c *echo.Context, message string) error {
	return ToastOnly(ctx, c, message, toast.VariantError)
}

// Redirect performs an HTTP redirect
func Redirect(c *echo.Context, url string) error {
	return c.Redirect(http.StatusFound, url)
}

// RedirectWithToast redirects with a toast notification using HX-Redirect header
func RedirectWithToast(ctx context.Context, c *echo.Context, url string, message string, variant toast.Variant) error {
	// Set HTMX redirect header
	c.Response().Header().Set("HX-Redirect", url)

	// Render toast with OOB swap (will appear after redirect)
	toastComponent := toast.Toast(toast.Props{
		Variant:       variant,
		Description:   message,
		Icon:          true,
		ShowIndicator: true,
	})

	oobToast := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := w.Write([]byte(`<div hx-swap-oob="afterbegin:#toast-container">`)); err != nil {
			return err
		}

		if err := toastComponent.Render(ctx, w); err != nil {
			return err
		}

		if _, err := w.Write([]byte(`</div>`)); err != nil {
			return err
		}

		return nil
	})

	return oobToast.Render(ctx, c.Response())
}

// HTMXRedirect sends an HTMX redirect header
func HTMXRedirect(c *echo.Context, url string) error {
	c.Response().Header().Set("HX-Redirect", url)
	return c.NoContent(http.StatusOK)
}

// JSON

// WriteJSON writes a JSON response with custom status code
func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

// --- Convenience Success Functions ---

func WriteCreated(c *echo.Context, data any) error {
	return WriteJSON(c.Response(), http.StatusCreated, data)
}

func WriteOK(c *echo.Context, data any) error {
	return WriteJSON(c.Response(), http.StatusOK, data)
}

func WriteNoContent(c *echo.Context) error {
	return WriteJSON(c.Response(), http.StatusNoContent, nil)
}

// --- Convenience Error Functions ---

func WriteBadRequest(c echo.Context, message, code string) error {
	return WriteErr(c, BadRequest(message, code))
}

func WriteNotFound(c echo.Context, message, code string) error {
	return WriteErr(c, NotFound(message, code))
}

func WriteUnauthorized(c echo.Context, message, code string) error {
	return WriteErr(c, Unauthorized(message, code))
}

func WriteForbidden(c echo.Context, message, code string) error {
	return WriteErr(c, Forbidden(message, code))
}

func WriteConflict(c echo.Context, message, code string) error {
	return WriteErr(c, Conflict(message, code))
}

// --- Core Error Handler ---

// WriteErr writes an error response.
// It accepts an optional message argument to override the default error message.
func WriteErr(c echo.Context, err error, messages ...string) error {
	var appErr *AppError

	// Type assert to your custom AppError
	if e, ok := err.(*AppError); ok {
		// Create a shallow copy to avoid mutating the original error object
		shallow := *e
		appErr = &shallow
	} else {
		// Use your Internal helper for unknown errors
		appErr = Internal("An unexpected error occurred", err)
	}

	if len(messages) > 0 {
		appErr.Message = messages[0]
	}

	res := c.Response()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(appErr.StatusCode)

	return json.NewEncoder(res).Encode(appErr)
}

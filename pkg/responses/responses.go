// pkg/responses/responses.go
package responses

import (
	"context"
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
	return c.NoContent(http.StatusTemporaryRedirect)
}

// JSON
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

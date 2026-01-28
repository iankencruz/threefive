// pkg/responses/responses.go
package responses

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/iankencruz/threefive/components/toast"
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

// RenderWithToast renders a component with a toast notification
func RenderWithToast(ctx context.Context, c *echo.Context, component templ.Component, message string, variant toast.Variant) error {
	// id main component
	if err := component.Render(ctx, c.Response()); err != nil {
		return err
	}

	// Render toast (out-of-band swap)
	toastComponent := toast.Toast(toast.Props{
		Variant:       variant,
		Description:   message,
		Icon:          true,
		ShowIndicator: true,
	})

	return toastComponent.Render(ctx, c.Response())
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

// ToastOnly renders only a toast notification
func ToastOnly(ctx context.Context, c *echo.Context, message string, variant toast.Variant) error {
	toastComponent := toast.Toast(toast.Props{
		Variant:       variant,
		Description:   message,
		Icon:          true,
		ShowIndicator: true,
	})

	return toastComponent.Render(ctx, c.Response())
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

// RedirectWithToast redirects with a toast notification
func RedirectWithToast(ctx context.Context, c *echo.Context, url string, message string, variant toast.Variant) error {
	c.Response().Header().Set("HX-Redirect", url)

	toastComponent := toast.Toast(toast.Props{
		Variant:     variant,
		Description: message,
	})

	return toastComponent.Render(ctx, c.Response())
}

// HTMXRedirect sends an HTMX redirect header
func HTMXRedirect(c *echo.Context, url string) error {
	c.Response().Header().Set("HX-Redirect", url)
	return c.NoContent(http.StatusOK)
}

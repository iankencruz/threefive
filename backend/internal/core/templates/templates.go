package templates

import (
	"net/http"

	"github.com/a-h/templ"
)

// RenderTempl renders a full HTML page using a layout and component.

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = component.Render(r.Context(), w)
}

package templates

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/iankencruz/threefive/internal/contextkeys"
)

// RenderTempl renders a full HTML page using a layout and component.

func Render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextkeys.CurrentPath, r.URL.Path)

	// if r.SessionManager != nil && r.UserModel != nil {
	// 	if userID, err := r.SessionManager.GetUserID(req); err == nil {
	// 		if user, err := r.UserModel.GetByID(ctx, r.DB, userID); err == nil {
	// 			ctx = context.WithValue(ctx, contextkeys.User, user)
	// 		}
	// 	}
	// }

	if err := page.Render(ctx, w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

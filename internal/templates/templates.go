package templates

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

type contextKey string

const (
	ctxUserKey        contextKey = "user"
	ctxSettingsKey    contextKey = "settings"
	ctxCurrentPathKey contextKey = "currentPath"
)

// RenderTempl renders a full HTML page using a layout and component.
func RenderTempl(
	w http.ResponseWriter,
	r *http.Request,
	layout func(title string, content templ.Component) templ.Component,
	title string,
	content templ.Component,
) {
	ctx := r.Context()

	// // Inject user (to be implemented later)
	// if userID, _ := GetSession(r); userID > 0 {
	//     if user, err := app.UserModel.GetUserByID(userID); err == nil {
	//         ctx = context.WithValue(ctx, ctxUserKey, user)
	//     }
	// }

	// Inject global settings if needed
	// if app.SettingsModel != nil {
	//     if settings, err := app.SettingsModel.GetAll(); err == nil {
	//         ctx = context.WithValue(ctx, ctxSettingsKey, settings)
	//     }
	// }

	// Inject current path
	ctx = context.WithValue(ctx, ctxCurrentPathKey, r.URL.Path)

	// Compose layout and render
	err := layout(title, content).Render(ctx, w)
	if err != nil {
		log.Printf("❌ RenderTempl error: %v", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

package server

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/iankencruz/threefive"
	"github.com/iankencruz/threefive/internal/handler"
	mw "github.com/iankencruz/threefive/internal/middleware"
	"github.com/labstack/echo/v5"
	"golang.org/x/time/rate"
)

func (s *Server) RegisterRoutes() {
	s.Log.Info("Registering routes...")
	s.Echo.Static("/assets", "assets")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(s.AuthService, s.SessionManager, s.Log)
	adminHandler := handler.NewAdminHandler(s.Log, s.MediaService)
	systemConfigHandler := handler.NewSystemConfigHandler(s.Log, s.SystemConfigService)
	mediaHandler := handler.NewMediaHandler(s.MediaService, s.Log)
	pageHandler := handler.NewPageHandler(s.Log, s.PageService, s.ProjectService, s.MediaService, s.SeoService)
	projectHandler := handler.NewProjectHandler(s.Log, s.ProjectService, s.TagService, s.MediaService, s.SeoService)
	tagHandler := handler.NewTagHandler(s.Log, s.TagService)
	contactHandler := handler.NewContactHandler(s.Log, s.ContactService)

	// Global Middleware
	s.Echo.Use(s.SessionMiddleware.Session)

	// 1. PUBLIC & AUTH ROUTES
	s.Echo.GET("/health", s.healthCheckHandler)

	loginLimiter := mw.RateLimit(rate.Limit(5.0/60), 5)
	s.Echo.GET("/login", authHandler.ShowLoginPage)
	s.Echo.POST("/login", authHandler.HandleLogin, loginLimiter)
	s.Echo.POST("/logout", authHandler.HandleLogout)

	s.Echo.GET("/", pageHandler.ShowPublicHome)
	s.Echo.GET("/about", pageHandler.ShowPublicAbout)
	s.Echo.GET("/projects", projectHandler.ShowPublicProjectsList)
	s.Echo.GET("/projects/:slug", projectHandler.ShowPublicProject)

	contactLimiter := mw.RateLimit(rate.Limit(10.0/3600), 3)
	s.Echo.GET("/contact", pageHandler.ShowPublicContact)
	s.Echo.POST("/contact", contactHandler.HandleSubmit, contactLimiter)

	// 2. PROTECTED API ROUTES (/api/admin/...)
	// This matches what your SvelteKit fetch calls are hitting
	api := s.Echo.Group("/api/admin")
	api.Use(s.SessionMiddleware.RequireAuth)

	api.GET("/dashboard", adminHandler.ShowDashboard)

	configGroup := api.Group("/system-config")
	configGroup.GET("", systemConfigHandler.ListSystemConfig)

	mediaGroup := api.Group("/media")
	mediaGroup.GET("", mediaHandler.ShowMediaList)
	mediaGroup.GET("/selector", mediaHandler.ShowMediaSelector)
	mediaGroup.POST("/upload", mediaHandler.UploadMedia)
	mediaGroup.GET("/:id/detail", mediaHandler.GetMediaDetail)
	mediaGroup.PUT("/:id", mediaHandler.UpdateMedia)
	mediaGroup.DELETE("/:id", mediaHandler.DeleteMedia)

	pagesGroup := api.Group("/pages")
	pagesGroup.GET("", pageHandler.ListPages)
	pagesGroup.GET("/:slug", pageHandler.ShowEditPage)
	pagesGroup.PUT("/:slug", pageHandler.UpdatePage)

	projectsGroup := api.Group("/projects")
	projectsGroup.GET("", projectHandler.ShowProjectsList)
	projectsGroup.POST("", projectHandler.CreateProject)
	projectsGroup.GET("/create-modal", projectHandler.ShowCreateModal)
	projectsGroup.GET("/gallery-selector", projectHandler.ShowGallerySelector)
	projectsGroup.GET("/:slug", projectHandler.ShowEditPage)
	projectsGroup.DELETE("/:slug", projectHandler.DeleteProject)
	projectsGroup.PUT("/:slug", projectHandler.UpdateProject)
	projectsGroup.PUT("/:slug/publish", projectHandler.PublishProject)
	projectsGroup.PUT("/:slug/unpublish", projectHandler.UnpublishProject)

	tagsGroup := api.Group("/tags")
	tagsGroup.GET("", tagHandler.ShowTagsList)
	tagsGroup.POST("", tagHandler.CreateTag)
	tagsGroup.GET("/create-modal", tagHandler.ShowCreateModal)
	tagsGroup.GET("/unused", tagHandler.ShowUnusedTags)
	tagsGroup.DELETE("/bulk/unused", tagHandler.DeleteUnusedTags)
	tagsGroup.GET("/:slug", tagHandler.ShowEditPage)
	tagsGroup.PUT("/:slug", tagHandler.UpdateTag)
	tagsGroup.DELETE("/:slug", tagHandler.DeleteTag)

	// 3. ADMIN SPA SERVING (/admin/...)
	// Register this LAST so it doesn't swallow other routes
	adminFS, err := fs.Sub(threefive.AdminAssets, "build")
	if err != nil {
		s.Log.Error("failed to create sub filesystem", "error", err)
	}

	adminUI := s.Echo.Group("/admin")
	adminUI.Use(s.SessionMiddleware.RequireAuth)

	fileServer := http.FileServer(http.FS(adminFS))

	adminUI.GET("/*", func(c *echo.Context) error {
		path := c.Request().URL.Path

		// Handle base path
		if path == "/admin" || path == "/admin/" {
			return c.FileFS("index.html", adminFS)
		}

		filePath := strings.TrimPrefix(path, "/admin/")

		// Check if file exists in the embedded FS (e.g., _app/immutable/...)
		f, err := adminFS.Open(filePath)
		if err != nil {
			// If file doesn't exist, serve index.html for SPA routing
			return c.FileFS("index.html", adminFS)
		}
		f.Close()

		return echo.WrapHandler(http.StripPrefix("/admin", fileServer))(c)
	})

	s.Log.Info("routes registered successfully")
}

func (s *Server) healthCheckHandler(c *echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

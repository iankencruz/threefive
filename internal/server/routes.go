package server

import (
	"github.com/iankencruz/threefive/internal/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	s.Log.Info("Registering routes...")
	s.Echo.Static("/assets", "assets")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(s.AuthService, s.SessionManager, s.Log)
	adminHandler := handler.NewAdminHandler(s.Log, s.MediaService)
	mediaHandler := handler.NewMediaHandler(s.MediaService, s.Log)
	pageHandler := handler.NewPageHandler(s.Log, s.PageService, s.ProjectService)
	projectHandler := handler.NewProjectHandler(s.Log, s.ProjectService, s.TagService, s.MediaService)
	tagHandler := handler.NewTagHandler(s.Log, s.TagService)

	// Static assets

	// Session Middleware
	s.Echo.Use(s.SessionMiddleware.Session)

	s.Echo.Pre(middleware.RemoveTrailingSlash())

	// handlers
	s.Echo.GET("/health", s.healthCheckHandler)

	// Public routes (no auth required)
	s.Echo.GET("/login", authHandler.ShowLoginPage)
	s.Echo.POST("/login", authHandler.HandleLogin)
	s.Echo.POST("/logout", authHandler.HandleLogout)

	s.Echo.GET("/", pageHandler.ShowPublicHome)
	s.Echo.GET("/about", pageHandler.ShowPublicAbout)
	s.Echo.GET("/projects", projectHandler.ShowPublicProjectsList)
	s.Echo.GET("/projects/:slug", projectHandler.ShowPublicProject)
	// s.Echo.GET("/contact", pageHandler.ShowPublicContact)

	// *********
	// admin routes (require authentication)
	// *********
	admin := s.Echo.Group("/admin")
	admin.Use(s.SessionMiddleware.RequireAuth)
	// redirect admin to dashboard
	admin.GET("", func(c *echo.Context) error {
		return c.Redirect(302, "/admin/dashboard")
	})

	// Dashboard Handler
	admin.GET("/dashboard", adminHandler.ShowDashboard)

	// Media Management
	media := admin.Group("/media")
	media.GET("", mediaHandler.ShowMediaList)
	media.GET("/selector", mediaHandler.ShowMediaSelector)
	media.POST("/upload", mediaHandler.UploadMedia)
	media.GET("/:id/detail", mediaHandler.GetMediaDetail)
	media.PUT("/:id", mediaHandler.UpdateMedia)
	media.DELETE("/:id", mediaHandler.DeleteMedia)

	// Page management (admin only)
	pages := admin.Group("/pages")

	pages.GET("", pageHandler.ListPages) // List all 3 pages
	pages.GET("/:slug", pageHandler.ShowEditPage)
	pages.PUT("/:slug", pageHandler.UpdatePage)

	// Project Management
	projects := admin.Group("/projects")

	projects.GET("", projectHandler.ShowProjectsList)
	projects.POST("", projectHandler.CreateProject)
	projects.GET("/create-modal", projectHandler.ShowCreateModal)
	projects.GET("/gallery-selector", projectHandler.ShowGallerySelector)
	projects.GET("/:slug", projectHandler.ShowEditPage)
	projects.DELETE("/:slug", projectHandler.DeleteProject)
	projects.PUT("/:slug", projectHandler.UpdateProject)
	projects.PUT("/:slug/publish", projectHandler.PublishProject)
	projects.PUT("/:slug/unpublish", projectHandler.UnpublishProject)

	// Tag Management
	tags := admin.Group("/tags")

	tags.GET("", tagHandler.ShowTagsList)
	tags.POST("", tagHandler.CreateTag)
	tags.GET("/create-modal", tagHandler.ShowCreateModal)
	tags.GET("/unused", tagHandler.ShowUnusedTags)
	tags.DELETE("/bulk/unused", tagHandler.DeleteUnusedTags)
	tags.GET("/:slug", tagHandler.ShowEditPage)
	tags.PUT("/:slug", tagHandler.UpdateTag)
	tags.DELETE("/:slug", tagHandler.DeleteTag)

	s.Log.Info("routes registered successfully")
}

// Simple handler methods on the server for basic routes
func (s *Server) healthCheckHandler(c *echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

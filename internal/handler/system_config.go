package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/iankencruz/threefive/internal/services"
	"github.com/labstack/echo/v5"
)

type SystemConfigHandler struct {
	logger              *slog.Logger
	SystemConfigService *services.SystemConfigService
}

func NewSystemConfigHandler(logger *slog.Logger, systemConfigService *services.SystemConfigService) *SystemConfigHandler {
	return &SystemConfigHandler{
		logger:              logger,
		SystemConfigService: systemConfigService,
	}
}

func (h *SystemConfigHandler) ListSystemConfig(c *echo.Context) error {
	h.logger.Debug("Loading System Config")

	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err != nil && parsed > 0 {
			page = parsed
		}
	}

	limit := int32(20)
	offset := int32((page - 1) * int(limit))

	configs, err := h.SystemConfigService.ListSystemConfig(c.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list configs", "error", err)
		return c.String(500, "Failed to load configs")
	}

	fmt.Printf("Configs: %v", configs)

	return c.JSON(http.StatusOK, configs)
}

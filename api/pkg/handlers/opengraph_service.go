package handlers

import (
	"context"
	"net/http"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/routes"
	"github.com/labstack/echo/v4"
)

type OpenGraphService interface {
	OpenGraphEditor(ctx context.Context, params routes.OpenGraphParams) (string, error)
	GetMetadata(ctx context.Context, params routes.GetMetadataParams) (routes.Metadata, error)
}

// OpenGraph - Data
// (GET /opengraph)
func (svc *Service) OpenGraph(c echo.Context, params routes.OpenGraphParams) error {

	html, err := svc.Services.OpenGraphSvc.OpenGraphEditor(c.Request().Context(), params)
	if err != nil {
		svc.logger.Error("Failed to get OpenGraph data:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get OpenGraph data")
	}

	return c.HTML(http.StatusOK, html)
}

func (svc *Service) GetMetadata(c echo.Context, params routes.GetMetadataParams) error {

	metadata, err := svc.Services.OpenGraphSvc.GetMetadata(c.Request().Context(), params)
	if err != nil {
		svc.logger.Error("Failed to get OpenGraph data:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get OpenGraph data")
	}

	return c.JSON(http.StatusOK, metadata)
}

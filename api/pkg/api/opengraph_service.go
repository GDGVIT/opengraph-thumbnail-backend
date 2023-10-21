package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OpenGraphService interface {
	OpenGraphEditor(ctx context.Context, params OpenGraphParams) (string, error)
}

// OpenGraph - Data
// (GET /opengraph)
func (svc *Service) OpenGraph(c echo.Context, params OpenGraphParams) error {

	html, err := svc.Services.OpenGraphSvc.OpenGraphEditor(c.Request().Context(), params)
	if err != nil {
		svc.logger.Error("Failed to get OpenGraph data:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get OpenGraph data")
	}

	return c.HTML(http.StatusOK, html)
}

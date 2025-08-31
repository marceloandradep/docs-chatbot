package ingest

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	ingestor *Ingestor
}

func NewHandler(ingestor *Ingestor) *Handler {
	return &Handler{
		ingestor: ingestor,
	}
}

func (h *Handler) Handle(c echo.Context) error {
	type req struct {
		Paths []string `json:"paths"`
		DocID string   `json:"doc_id"`
	}

	var body req
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload format")
	}

	if len(body.Paths) == 0 || body.DocID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "doc_id and paths required")
	}

	count, err := h.ingestor.IngestPaths(context.Background(), body.DocID, body.Paths)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{"ingested_chunks": count})
}

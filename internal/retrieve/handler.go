package retrieve

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Handler struct {
	retriever *Retriever
}

func NewHandler(retriever *Retriever) *Handler {
	return &Handler{
		retriever: retriever,
	}
}

func (r *Handler) Handle(c echo.Context) error {
	type req struct {
		Query     string `json:"query"`
		DocID     string `json:"doc_id"`
		TopK      int    `json:"top_k"`
		MaxTokens int    `json:"max_tokens"`
	}

	var body req
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload format")
	}

	if strings.TrimSpace(body.Query) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query required")
	}

	if body.TopK == 0 {
		body.TopK = 6
	}

	if body.MaxTokens == 0 {
		body.MaxTokens = 400
	}

	answer, cites, err := r.retriever.Answer(context.Background(), AnswerParams{
		Query:     body.Query,
		DocID:     body.DocID,
		TopK:      body.TopK,
		MaxTokens: body.MaxTokens,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"answer":    answer,
		"citations": cites,
	})
}

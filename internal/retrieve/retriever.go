package retrieve

import (
	"chatbot/internal/commons"
	"chatbot/internal/ingest"
	"context"
	"database/sql"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type AnswerParams struct {
	Query     string
	DocID     string
	TopK      int
	MaxTokens int
}

type Retrieved struct {
	Content    string
	SourcePath string
	Idx        int
}

type Retriever struct {
	db        *sql.DB
	oaiCli    *openai.Client
	ingestor  *ingest.Ingestor
	chatModel string
}

func NewRetriever(db *sql.DB, oaiCli *openai.Client, ingestor *ingest.Ingestor, chatModel string) *Retriever {
	return &Retriever{
		db:        db,
		oaiCli:    oaiCli,
		ingestor:  ingestor,
		chatModel: chatModel,
	}
}

func (r *Retriever) Answer(ctx context.Context, p AnswerParams) (string, []map[string]any, error) {
	qEmb, err := r.ingestor.EmbedOne(ctx, p.Query)
	if err != nil {
		return "", nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT content, source_path, chunk_index
		FROM chunks
		WHERE ($1 = '' OR doc_id = $1)
		ORDER BY embedding <=> $2::vector
		LIMIT $3
	`, p.DocID, commons.FsToString(qEmb), p.TopK)
	defer rows.Close()
	if err != nil {
		return "", nil, err
	}

	var ctxs []Retrieved
	for rows.Next() {
		var r Retrieved
		if err := rows.Scan(&r.Content, &r.SourcePath, &r.Idx); err != nil {
			return "", nil, err
		}
		ctxs = append(ctxs, r)
	}

	system := `You answer strictly from the given context. If unsure, say you don't know. Return concise answers and list sources.`
	ctxBlock := ""
	for i, r := range ctxs {
		ctxBlock += fmt.Sprintf("\n[CTX %d] (%s#%d)\n%s\n", i+1, r.SourcePath, r.Idx, r.Content)
	}
	user := fmt.Sprintf("Question: %s\n\nContext:\n%s\n\nInstructions: cite sources as (source#chunk).", p.Query, ctxBlock)

	resp, err := r.oaiCli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:               r.chatModel,
		MaxCompletionTokens: p.MaxTokens,
		Temperature:         0.1,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	})
	if err != nil {
		return "", nil, err
	}

	cites := make([]map[string]any, len(ctxs))
	for i, r := range ctxs {
		snippet := r.Content
		if len(snippet) > 220 {
			snippet = snippet[:220] + "..."
		}
		cites[i] = map[string]any{
			"source_path": r.SourcePath,
			"chunk_index": r.Idx,
			"preview":     snippet,
		}
	}

	return resp.Choices[0].Message.Content, cites, nil
}

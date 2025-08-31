package ingest

import (
	"chatbot/internal/store"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/sashabaranov/go-openai"
	"os"
	"path/filepath"
	"strings"
)

const (
	approximateChunkSize = 1200
	approximateOverlap   = 200
)

type Ingestor struct {
	db       *store.Store
	cli      *openai.Client
	embModel string
}

func NewIngestor(db *store.Store, cli *openai.Client, embModel string) *Ingestor {
	return &Ingestor{
		db:       db,
		cli:      cli,
		embModel: embModel,
	}
}

func (uc *Ingestor) IngestPaths(ctx context.Context, docID string, paths []string) (int, error) {
	var total int
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return total, err
		}
		text := normalize(string(data), filepath.Ext(p))
		chunks := chunk(text, approximateChunkSize, approximateOverlap)
		embs, err := uc.embed(ctx, chunks)
		if err != nil {
			return total, err
		}
		for i, v := range embs {
			if err := uc.db.InsertChunk(ctx, docID, p, i, chunks[i], v); err != nil {
				log.Error(err)
				return total, err
			}
			total++
		}
	}
	return total, nil
}

func (uc *Ingestor) EmbedOne(ctx context.Context, text string) ([]float32, error) {
	embs, err := uc.embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return embs[0], nil
}

func (uc *Ingestor) embed(ctx context.Context, chunks []string) ([][]float32, error) {
	resp, err := uc.cli.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(uc.embModel),
		Input: chunks,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	out := make([][]float32, len(resp.Data))
	for i, d := range resp.Data {
		out[i] = d.Embedding
	}
	return out, nil
}

// placeholder for more sophisticated normalization routine
func normalize(s, ext string) string {
	if strings.EqualFold(ext, ".md") {
		return s
	}
	return s
}

func chunk(s string, size, overlap int) []string {
	var out []string
	for i := 0; i < len(s); i += size - overlap {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
		if end == len(s) {
			break
		}
	}
	return out
}

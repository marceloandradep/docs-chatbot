package registry

import (
	"chatbot/internal/env"
	"chatbot/internal/ingest"
)

const (
	defaultModel = "text-embedding-3-large"
)

func (l *Locator) NewIngestHandler() *ingest.Handler {
	return ingest.NewHandler(l.NewIngestor())
}

func (l *Locator) NewIngestor() *ingest.Ingestor {
	embModel := env.GetOrDefault(env.EmbeddingModel, defaultModel)
	return ingest.NewIngestor(l.NewStore(), l.oaiCli, embModel)
}

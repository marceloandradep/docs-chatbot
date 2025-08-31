package env

import "os"

const (
	OpenAPIKey         = "OPENAI_API_KEY"
	DataBaseURL        = "DATABASE_URL"
	EmbeddingDimension = "EMBED_DIMS"
	EmbeddingModel     = "EMBEDDING_MODEL"
	ChatModel          = "CHAT_MODEL"
	AppPort            = "PORT"
)

func GetOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

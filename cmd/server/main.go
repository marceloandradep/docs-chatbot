package main

import (
	"chatbot/internal/env"
	"chatbot/internal/registry"
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
)

func main() {
	cfg := mustLoadConfig()

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	oaiClient := openai.NewClient(cfg.OpenAIKey)
	locator := registry.NewLocator(db, oaiClient)

	e := echo.New()

	e.POST("/ingest", locator.NewIngestHandler().Handle)
	e.POST("/ask", locator.NewRetrievalHandler().Handle)

	log.Fatal(e.Start(":" + cfg.Port))
}

type Config struct {
	DbURL     string
	OpenAIKey string
	EmbedDims int
	Port      string
}

func mustLoadConfig() Config {
	key := os.Getenv(env.OpenAPIKey)
	if key == "" {
		log.Fatal("OpenAPI API key is required")
	}

	db := env.GetOrDefault(env.DataBaseURL, "postgres://postgres:postgres@localhost:5432/vecdb?sslmode=disable")
	dims, _ := strconv.Atoi(env.GetOrDefault(env.EmbeddingDimension, "3072"))
	port := env.GetOrDefault(env.AppPort, "8080")

	return Config{
		DbURL:     db,
		OpenAIKey: key,
		EmbedDims: dims,
		Port:      port,
	}
}

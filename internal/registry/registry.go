package registry

import (
	"database/sql"
	"github.com/sashabaranov/go-openai"
)

type Locator struct {
	db     *sql.DB
	oaiCli *openai.Client
}

func NewLocator(db *sql.DB, oaiCli *openai.Client) *Locator {
	return &Locator{
		db:     db,
		oaiCli: oaiCli,
	}
}

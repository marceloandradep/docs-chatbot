package store

import (
	"chatbot/internal/commons"
	"context"
	"database/sql"
)

type Store struct {
	db  *sql.DB
	dim int
}

func NewStore(db *sql.DB, dim int) *Store {
	return &Store{
		db:  db,
		dim: dim,
	}
}

func (s *Store) InsertChunk(ctx context.Context, docID, path string, idx int, content string, emb []float32) error {
	str := commons.FsToString(emb)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chunks (doc_id, source_path, chunk_index, content, embedding)
		VALUES ($1, $2, $3, $4, $5::vector)
	`, docID, path, idx, content, str)
	return err
}

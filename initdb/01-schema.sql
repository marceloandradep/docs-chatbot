CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS chunks (
    id          BIGSERIAL PRIMARY KEY,
    doc_id      TEXT NOT NULL,
    source_path TEXT NOT NULL,
    chunk_index INT  NOT NULL,
    content     TEXT NOT NULL,
    embedding   VECTOR(3072) NOT NULL,
    metadata    JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_chunks_embedding
    ON chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_chunks_docid ON chunks (doc_id);
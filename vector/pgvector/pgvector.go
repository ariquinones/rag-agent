package pgvector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rag-agent/vector"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

const (
	ENSURE_DB_EXTENSION = "CREATE EXTENSION IF NOT EXISTS vector;"
)

type Options struct {
	DSN          string // connection string to postgres
	EmbeddingDim int    // must match the embeddor
}

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, opts Options) (*Store, error) {
	if opts.DSN == "" {
		return nil, errors.New("pgvector dsn is required")
	}
	if opts.EmbeddingDim <= 0 {
		return nil, errors.New("pgvector embedding dim must be greater than 0")
	}
	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	if err := ensureVectorExtension(ctx, opts.DSN); err != nil {
		return nil, fmt.Errorf("install vector extension: %w", err)
	}
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	s := &Store{
		pool: pool,
	}
	if err := s.migrate(ctx, opts.EmbeddingDim); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func ensureVectorExtension(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, ENSURE_DB_EXTENSION)
	return err
}

func (s *Store) migrate(ctx context.Context, dim int) error {
	statements := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			embedding vector(%d) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now())`, dim),
		`CREATE INDEX IF NOT EXISTS documents_embedding_idx
			  ON documents USING hnsw (embedding vector_cosine_ops)`,
	}
	for _, query := range statements {
		if _, err := s.pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("exec %q: %w", query, err)
		}
	}
	return nil
}

func (s *Store) Upsert(ctx context.Context, docs []vector.Document) error {
	if len(docs) == 0 {
		return nil
	}
	transaction, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	const stmt = `
		INSERT INTO documents (id, content, metadata, embedding)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			embedding = EXCLUDED.embedding`
	for _, doc := range docs {
		meta, err := marshalMetadata(doc.Metadata)
		if err != nil {
			return fmt.Errorf("metadata for: %s %w", doc.ID, err)
		}

		if _, err := transaction.Exec(ctx, stmt, doc.ID, doc.Content, meta, pgvector.NewVector(doc.Embedding)); err != nil {
			return fmt.Errorf("upsert: %s %w", doc.ID, err)
		}
	}
	return transaction.Commit(ctx)
}

func marshalMetadata(m map[string]string) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

func unmarshalMetadat(bytes []byte, destination *map[string]string) error {
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, destination)
}

func (s *Store) Query(ctx context.Context, embedding []float32, topK int) ([]vector.Result, error) {
	if topK <= 0 {
		return nil, nil
	}
	const stmt = `
		select id, content, metadata, embedding <=> $1 as distance
		from documents
		order by embedding <=> $1
		limit $2
	`
	rows, err := s.pool.Query(ctx, stmt, pgvector.NewVector(embedding), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []vector.Result{}
	for rows.Next() {
		var (
			r        vector.Result
			metaRaw  []byte
			distance float64
		)
		if err := rows.Scan(&r.ID, &r.Content, &metaRaw, &distance); err != nil {
			return nil, err
		}
		if err := unmarshalMetadat(metaRaw, &r.Metadata); err != nil {
			return nil, fmt.Errorf("metadata for: %s %w", r.ID, err)
		}
		// set the score to a float32 so the caller has a better way of sorting on their end
		r.Score = float32(1 - distance)
		results = append(results, r)
	}
	return results, rows.Err()
}

func (s *Store) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.pool.Exec(ctx, `delete from documents where id = ANY($1)`, ids)
	return err
}

func (s *Store) DeleteBySource(ctx context.Context, source string) error {
	if source == "" {
		return nil
	}
	_, err := s.pool.Exec(ctx, `delete from documents where metadata->>'source' = $1`, source)
	return err
}

func (s *Store) Close() error {
	s.pool.Close()
	return nil
}

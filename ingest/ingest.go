package ingest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rag-agent/llm"
	"rag-agent/vector"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_CHUNK_SIZE    = 1000 // 1k bytes
	DEFAULT_CHUNK_OVERLAP = 100  // 100 bytes
)

type Options struct {
	SourceDir    string
	ProcessedDir string
	ChunkSize    int
	ChunkOverlap int
}

func processFile(ctx context.Context, path string, opts Options, embedder llm.Embedder, store vector.Store) error {
	if !supportedFormat(path) {
		return fmt.Errorf("unsupported file format: %s", filepath.Ext(path))
	}
	// read file
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	_, err = processContent(ctx, path, raw, opts, embedder, store)
	return err
}

// processContent - implements an ingest pipeline to take a document from disk and puts them in vector storebackend
// 1 read file content
// 2 chunk the file
// 3 embed the chunks
// 4 delete any prior chunks from file
// 5 upsert into the db store
func processContent(ctx context.Context, sourceFilePath string, content []byte, opts Options, embedder llm.Embedder, store vector.Store) (int, error) {
	if embedder == nil {
		return 0, errors.New("embedder is required")
	}
	if store == nil {
		return 0, errors.New("vector store is required")
	}
	base := filepath.Base(sourceFilePath)
	if !supportedFormat(base) {
		return 0, fmt.Errorf("unsupported file format: %s", filepath.Ext(base))
	}
	size := opts.ChunkSize
	if size <= 0 {
		size = DEFAULT_CHUNK_SIZE
	}
	overlap := opts.ChunkOverlap
	if overlap <= 0 {
		overlap = DEFAULT_CHUNK_OVERLAP
	}
	text := strings.TrimSpace(string(content))
	if text == "" {
		return 0, errors.New("file is empty")
	}
	chunks := chunk(text, size, overlap)
	if len(chunks) == 0 {
		return 0, errors.New("no chunks produced")
	}

	vectorEmbeddings, err := embedder.Embed(ctx, chunks)
	if err != nil {
		return 0, fmt.Errorf("emdedding chunks issue: %w", err)
	}
	if len(vectorEmbeddings) != len(chunks) {
		return 0, fmt.Errorf("emded: got %d vectors for %d chunks", len(vectorEmbeddings), len(chunks))
	}

	if err := store.DeleteBySource(ctx, sourceFilePath); err != nil {
		return 0, fmt.Errorf("clear previous chunks: %w", err)
	}
	ingestedAt := time.Now().UTC().Format(time.RFC3339)
	docs := make([]vector.Document, len(chunks))
	for i, chunk := range chunks {
		// Corrected: strings.ToValidUTF8 replaces invalid UTF-8 byte sequences.
		// Setting the second argument to "" will completely strip out the invalid bytes (like 0x92).
		// If you want a visual placeholder instead, you can pass "" or "?".
		sanitizedContent := strings.ToValidUTF8(chunk, "")
		docs[i] = vector.Document{
			ID:      fmt.Sprintf("%s#%d", base, i),
			Content: sanitizedContent,
			Metadata: map[string]string{
				"source":      base,
				"chunk_index": strconv.Itoa(i),
				"chunks":      strconv.Itoa(len(chunks)),
				"ingested_at": ingestedAt,
			},
			Embedding: vectorEmbeddings[i],
		}
	}
	err = store.Upsert(ctx, docs)
	if err != nil {
		return 0, err
	}
	return len(chunks), nil
}

func supportedFormat(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".md", ".markdown":
		return true
	}
	return false
}

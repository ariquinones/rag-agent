package app

import (
	"context"
	"log"
	"os"
	"rag-agent/chat"
	"rag-agent/config"
	"rag-agent/ingest"
	"rag-agent/llm"
	"rag-agent/rag"
	"rag-agent/vector"
	"rag-agent/vector/pgvector"
	"sync"
)

func Run(parent context.Context, cfg config.Config) error {
	logger := log.New(os.Stderr, "[rag-agent] ", log.LstdFlags)

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	client := llm.New(cfg)
	embedderClient := llm.NewEmbedder(cfg)

	store, err := openStore(ctx, cfg)
	if err != nil {
		logger.Printf("vector store disabled %v", err)
	}
	var wg sync.WaitGroup
	if store != nil {
		defer store.Close()
		logger.Printf("vector store ready")

		wg.Go(func() {
			opts := ingest.Options{
				SourceDir:    cfg.IngestDir,
				ProcessedDir: cfg.ProcessedDir,
			}
			err := ingest.Watch(ctx, opts, embedderClient, store, logger)
			if err != nil {
				logger.Printf("watcher stopped: %v", err)
			}
		})
		logger.Printf("watching %s for new documents", cfg.IngestDir)
	}

	// get rag retriver
	var retriever *rag.Retriever
	if store != nil {
		retriever = rag.New(embedderClient, store, rag.Options{
			TopK:     5, // set default here for now...
			Rewriter: rag.NewRewriter(client),
		})
	}

	replErr := chat.RunEPL(ctx, client, retriever, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
	cancel()
	wg.Wait()
	return replErr
}

func openStore(ctx context.Context, cfg config.Config) (vector.Store, error) {
	if cfg.DatabaseURL == "" {
		return nil, nil
	}
	s, err := pgvector.New(ctx, pgvector.Options{
		DSN:          cfg.DatabaseURL,
		EmbeddingDim: cfg.EmbeddingDim,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

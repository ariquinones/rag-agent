## Overview

A terminal RAG chat agent in Go (module `rag-agent`). It runs a REPL against any OpenAI-compatible chat API, and in the background watches a directory for documents, then chunks, embeds and stores them in Postgres/pgvector. Each chat turn retrieves the closest chunks and adds them to the prompt.

## Commands

```sh
docker compose up -d          # Postgres 18 + pgvector on :5432 (user/pass/db: rag)
go run ./cmd/rag              # start the REPL (reads .env via godotenv)
go build ./...
go vet ./...
go test ./...                 # no tests exist yet; single test: go test ./ingest -run TestName
```

Quit the REPL with `q`, `quit`, `exit` or `/exit`.

## Configuration

`config.Load()` reads env vars (or `.env`): `OPENAI_BASE_URL`, `OPENAI_API_KEY`, `OPENAI_MODEL` (default `gpt-4o-mini`), `SYSTEM_PROMPT_FILE` (e.g. `prompts/system-prompt.md`; loaded as the system message), `DATABASE_URL`, `EMBEDDING_DIM` (default 768), `EMBEDDING_BASE_URL` / `EMBEDDING_API_KEY` (fall back to the chat values), `EMBEDDING_MODEL` (default `nomic-embed-text`), `INGEST_DIR` (default `./documents`), `PROCESSED_DIR` (default `./documents/processed`).

Chat and embeddings can point at different servers (e.g. a hosted chat model plus a local Ollama for `nomic-embed-text`). `EMBEDDING_DIM` must match the embedding model's output width, because it sets the `vector(N)` column. Changing models with a different width means dropping the `documents` table.

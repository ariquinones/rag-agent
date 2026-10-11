package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	NOMIC_EMBED_TEXT      = 768
	DEFAULT_BASE_URL      = "https://api.openai.com/v1"
	DEFAULT_MODEL         = "gpt-4o-mini"
	DEFAULT_EMBED_MODEL   = "nomic-embed-text"
	DEFAULT_INGEST_DIR    = "./documents"
	DEFAULT_PROCESSED_DIR = "./documents/processed"
)

type Config struct {
	BaseURL          string // any open ai compatible server
	APIkey           string
	Model            string
	SystemPromptFile string
	DatabaseURL      string // db connection
	EmbeddingDim     int    // embedding dimensions ex 768 (width of colum for the embeddings in the db)
	EmbeddingBaseUrl string
	EmbeddingApiKey  string
	EmbeddingModel   string
	IngestDir        string // directory to watch
	ProcessedDir     string // directory of processed files
}

func Load() Config {
	_ = godotenv.Load()
	cfg := Config{
		BaseURL:          os.Getenv("OPENAI_BASE_URL"),
		APIkey:           os.Getenv("OPENAI_API_KEY"),
		Model:            os.Getenv("OPENAI_MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		EmbeddingDim:     atoiOr(os.Getenv("EMBEDDING_DIM"), NOMIC_EMBED_TEXT),
		EmbeddingBaseUrl: os.Getenv("EMBEDDING_BASE_URL"),
		EmbeddingApiKey:  os.Getenv("EMBEDDING_API_KEY"),
		EmbeddingModel:   os.Getenv("EMBEDDING_MODEL"),
		IngestDir:        os.Getenv("INGEST_DIR"),
		ProcessedDir:     os.Getenv("PROCESSED_DIR"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DEFAULT_BASE_URL
	}
	if cfg.Model == "" {
		cfg.Model = DEFAULT_MODEL
	}
	if cfg.EmbeddingBaseUrl == "" {
		cfg.EmbeddingBaseUrl = cfg.BaseURL
	}
	if cfg.EmbeddingApiKey == "" {
		cfg.EmbeddingApiKey = cfg.APIkey
	}

	if cfg.EmbeddingModel == "" {
		cfg.EmbeddingModel = DEFAULT_EMBED_MODEL
	}

	if cfg.IngestDir == "" {
		cfg.IngestDir = DEFAULT_INGEST_DIR
	}
	if cfg.ProcessedDir == "" {
		cfg.ProcessedDir = DEFAULT_PROCESSED_DIR
	}

	return cfg
}

func atoiOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

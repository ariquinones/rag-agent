package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rag-agent/app"
	"rag-agent/config"
	"syscall"
)

func main() {
	// set up app
	// set up config
	// set up llm client
	// set up read-eval-print loop - REPL
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, config.Load()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

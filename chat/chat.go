package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"rag-agent/llm"
	"rag-agent/rag"
	"strings"
	"sync"
	"time"
)

type Options struct {
	SystemPromptFile string
}

func RunEPL(ctx context.Context, client *llm.Client, retriever *rag.Retriever, options Options) error {
	input := bufio.NewScanner(os.Stdin)
	input.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	chatHistory, err := seedHistory(options.SystemPromptFile)
	if err != nil {
		return err
	}
	fmt.Println("Session starting. Type Q to quit")

	for {
		fmt.Print("\n> ")
		if !input.Scan() {
			if err := input.Err(); err != nil {
				return err
			}
			return nil
		}
		textInput := strings.TrimSpace(input.Text())
		if textInput == "" {
			continue
		}

		// Handle exit commands
		if strings.EqualFold(textInput, "q") || strings.EqualFold(textInput, "/exit") || strings.EqualFold(textInput, "exit") || strings.EqualFold(textInput, "quit") {
			fmt.Println("Goodbye...")
			return nil
		}

		chatHistory = append(chatHistory, llm.Message{Role: "user", Content: textInput})
		turn := chatHistory
		thinking := startLoadingSpinner("thinking...")
		var stopOnce sync.Once

		if retriever != nil {
			contextText, rtrErr := retriever.Retrieve(ctx, chatHistory)
			if rtrErr != nil {
				fmt.Fprintln(os.Stderr, "retrieval err: %w", rtrErr)
			} else if contextText != "" {
				// build a turn with the inline context
				turn = withInlineContext(chatHistory, contextText)
			}
		}

		reply, err := client.ChatStream(ctx, turn, func(s string) {
			stopOnce.Do(thinking.Stop)
			fmt.Print(s)
		})
		stopOnce.Do(thinking.Stop)
		fmt.Println()

		if err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
			chatHistory = chatHistory[:len(chatHistory)-1]
			continue
		}

		chatHistory = append(chatHistory, reply)
	}

}

type loadingSpinner struct {
	stop chan struct{}
	done chan struct{}
	once sync.Once
}

func (s *loadingSpinner) Stop() {
	s.once.Do(func() { close(s.stop) })
	<-s.done
}

func withInlineContext(history []llm.Message, contextText string) []llm.Message {
	if len(history) == 0 || contextText == "" {
		return history
	}
	last := history[len(history)-1]
	if last.Role != "user" {
		return history
	}
	out := make([]llm.Message, len(history))
	copy(out, history)
	out[len(out)-1] = llm.Message{
		Role:    "user",
		Content: contextText + "\n\n --- Question --- \n\n" + last.Content,
	}
	return out
}
func startLoadingSpinner(label string) *loadingSpinner {
	s := &loadingSpinner{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		defer close(s.done)
		frames := []string{"|", "/", "-", "\\"}
		t := time.NewTicker(80 * time.Millisecond)
		defer t.Stop()
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Print("\r\033[K")
				return
			case <-t.C:
				fmt.Print("\r%s %s", frames[i%len(frames)], label)
				i++
			}
		}
	}()
	return s
}

func seedHistory(path string) ([]llm.Message, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("seed history err: %w", err)
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return nil, nil
	}
	return []llm.Message{{Role: "system", Content: content}}, nil
}

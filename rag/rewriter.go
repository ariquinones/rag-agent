package rag

import (
	"context"
	"fmt"
	"rag-agent/llm"
	"strings"
)

const REWRITE_SYSTEM_PROMPT = `You rewrite the user's last message into a standalone search query.

Given the conversation, output a single search query that:
- Captures the topic and intent of the latest user message.
- Resolves pronouns and refrences using prior turns ("it", "they", "that one").
- Stay concise - keywords and short phrases, not full sentences.

If the latest user message already stands on its own with no references to prior turns, output it verbatim.

Output only the query. No preamble, no quotes, no explanation.`

type Rewriter struct {
	client *llm.Client
}

func NewRewriter(client *llm.Client) *Rewriter {
	return &Rewriter{
		client: client,
	}
}

func (r *Rewriter) Rewrite(ctx context.Context, history []llm.Message) (string, error) {
	lastMsg := lastUserMessage(history)
	if lastMsg == "" {
		return "", nil
	}
	if !hasAssistantTurn(history) {
		return lastMsg, nil
	}
	msgs := []llm.Message{
		{Role: "system", Content: REWRITE_SYSTEM_PROMPT},
		{Role: "user", Content: formatConvo(history)},
	}
	reply, err := r.client.ChatStream(ctx, msgs, nil)
	if err != nil {
		return "", fmt.Errorf("rewrite err: %w", err)
	}
	out := strings.TrimSpace(reply.Content)
	out = strings.Trim(out, `"'`)
	if out == "" {
		return lastMsg, nil
	}
	return out, nil
}

func hasAssistantTurn(history []llm.Message) bool {
	for _, msg := range history {
		if msg.Role == "assistant" {
			return true
		}
	}
	return false
}

func formatConvo(history []llm.Message) string {
	var sb strings.Builder
	sb.WriteString("Conversation so far:\n\n")
	for _, msg := range history {
		switch msg.Role {
		case "user":
			sb.WriteString("User: ")
		case "assistant":
			sb.WriteString("Assistant: ")
		default:
			continue // skip system
		}
		sb.WriteString(msg.Content)
		sb.WriteString("\n\n")
	}
	sb.WriteString("Rewrite the user's latest message as a standalone search query.")
	return sb.String()
}

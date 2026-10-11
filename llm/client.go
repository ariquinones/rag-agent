package llm

import (
	"context"
	"rag-agent/config"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Message struct {
	Role    string `json:"role"`    // system, user or assistant
	Content string `json:"content"` // payload of message
}

type Client struct {
	Cfg config.Config
	sdk openai.Client
}

func New(cfg config.Config) *Client {
	return newClient(cfg, cfg.BaseURL, cfg.APIkey)
}

func NewEmbedder(cfg config.Config) *Client {
	return newClient(cfg, cfg.EmbeddingBaseUrl, cfg.EmbeddingApiKey)
}

func newClient(cfg config.Config, baseUrl string, apiKey string) *Client {
	ops := []option.RequestOption{}

	if baseUrl != "" {
		ops = append(ops, option.WithBaseURL(baseUrl))
	}
	if apiKey != "" {
		ops = append(ops, option.WithAPIKey(apiKey))
	}
	return &Client{
		Cfg: cfg,
		sdk: openai.NewClient(ops...),
	}
}

func (c *Client) ChatStream(ctx context.Context, messages []Message, onChange func(string)) (Message, error) {
	stream := c.sdk.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    c.Cfg.Model,
		Messages: toSdkMessages(messages),
	})
	defer stream.Close()

	var content strings.Builder
	role := "assistant"
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta
		if delta.Role != "" {
			role = delta.Role
		}
		if delta.Content != "" {
			content.WriteString(delta.Content)
			if onChange != nil {
				onChange(delta.Content)
			}
		}
	}
	if err := stream.Err(); err != nil {
		return Message{}, err
	}
	return Message{Role: role, Content: content.String()}, nil
}

func toSdkMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			out = append(out, openai.SystemMessage(msg.Content))
		case "assistant":
			out = append(out, openai.AssistantMessage(msg.Content))
		default:
			out = append(out, openai.UserMessage(msg.Content))
		}
	}
	return out
}

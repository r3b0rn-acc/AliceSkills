package openai

import (
	"context"
	"errors"
	"github.com/openai/openai-go/v2"
)

func New(cfg Config) *Client {
	cfg.Normalize()
	opts := cfg.RequestOptions()
	return &Client{
		inner:   openai.NewClient(opts...),
		model:   cfg.Model,
		timeout: cfg.Timeout,
		system:  cfg.System,
	}
}

func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	if c.timeout > 0 {
		if _, has := ctx.Deadline(); !has {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.timeout)
			defer cancel()
		}
	}

	params := openai.ChatCompletionNewParams{
		Messages: toOpenAI(messages),
		Model:    c.model,
	}

	resp, err := c.inner.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("empty choices")
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	messages := make([]Message, 0, 2)
	if c.system != "" {
		messages = append(messages, Message{Role: "system", Content: c.system})
	}
	messages = append(messages, Message{Role: "user", Content: prompt})
	return c.Chat(ctx, messages)
}

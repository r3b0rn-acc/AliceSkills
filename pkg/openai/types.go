package openai

import (
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"os"
	"time"
)

type SystemPrompt string

type Message struct {
	Role    string // "system" | "user" | "assistant"
	Content string
}

type Config struct {
	APIKey     string
	BaseURL    string
	Model      openai.ChatModel
	Timeout    time.Duration
	MaxRetries int
	Debug      bool
	System     string
}

func (cfg *Config) Normalize() {
	if cfg.Model == "" {
		cfg.Model = openai.ChatModelGPT4o
	}

}

func (cfg *Config) RequestOptions() []option.RequestOption {
	opts := make([]option.RequestOption, 0, 6)

	key := cfg.APIKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key != "" {
		opts = append(opts, option.WithAPIKey(key))
	}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	if cfg.MaxRetries >= 0 {
		opts = append(opts, option.WithMaxRetries(cfg.MaxRetries))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, option.WithRequestTimeout(cfg.Timeout))
	}
	if cfg.Debug {
		opts = append(opts, option.WithDebugLog(nil))
	}
	return opts
}

type Client struct {
	inner   openai.Client
	model   openai.ChatModel
	timeout time.Duration
	system  string
}

package openai

import (
	"sync"
	"time"
)

var (
	onceDefault   sync.Once
	onceInit      sync.Once
	defaultClient *Client
	initCfg       *Config
)

func Init(c Config) {
	onceInit.Do(func() {
		cfg := c
		initCfg = &cfg
	})
}

func Default() Client {
	onceDefault.Do(func() {
		var cfg Config
		if initCfg != nil {
			cfg = *initCfg
		} else {
			cfg = Config{
				Model:      "",
				Timeout:    30 * time.Second,
				MaxRetries: 2,
				System:     string(VoiceHelperSystemPrompt),
			}
		}
		defaultClient = New(cfg)
	})
	return *defaultClient
}

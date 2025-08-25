package skills

import (
	"context"
	"encoding/json"
)

type Skill interface {
	Name() string
	Handle(ctx context.Context, raw json.RawMessage) (any, error)
}

type TextFn func(string) string

type BaseSkill struct {
	Pipe TextFn
}

func (bs BaseSkill) Handle(_ context.Context, raw json.RawMessage) (any, error) {
	return baseHandler(raw, bs.Pipe)
}

func baseHandler(raw json.RawMessage, pipe TextFn) (any, error) {
	var input Event
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, err
	}

	text := pipe(input.Request.OriginalUtterance)

	return &Response{
		Version: input.Version,
		Session: input.Session,
		Result: struct {
			Text       string `json:"text"`
			EndSession bool   `json:"end_session"`
		}{
			Text:       text,
			EndSession: false,
		},
	}, nil
}

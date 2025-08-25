package gptquestion

import (
	"AliceSkills/internal/skills"
	"AliceSkills/pkg/openai"
	"context"
)

type Skill struct{ skills.BaseSkill }

func (Skill) Name() string { return "gpt_question" }

func New() *Skill {
	return &Skill{BaseSkill: skills.BaseSkill{Pipe: gptResponseText}}
}

func gptResponseText(inputText string) string {
	client := openai.Default()
	reply, err := client.Chat(context.Background(), []openai.Message{
		{Role: "system", Content: "Отвечай коротко и по делу."},
		{Role: "user", Content: inputText},
	})
	if err != nil {
		panic(err)
	}
	return reply
}

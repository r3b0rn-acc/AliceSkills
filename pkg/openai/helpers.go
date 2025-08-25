package openai

import "github.com/openai/openai-go/v2"

const VoiceHelperSystemPrompt SystemPrompt = "Ты голосовой помощник. Отвечай понятно, собранно, только по делу, но в то же время кратко."

func toOpenAI(messages []Message) []openai.ChatCompletionMessageParamUnion {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "system":
			out = append(out, openai.SystemMessage(m.Content))
		case "assistant":
			out = append(out, openai.AssistantMessage(m.Content))
		default: // "user"
			out = append(out, openai.UserMessage(m.Content))
		}
	}
	return out
}

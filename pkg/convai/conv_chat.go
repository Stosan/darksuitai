package convai

import (
	ant "github.com/Stosan/darksuitai/internal/llms/anthropic"
	oai "github.com/Stosan/darksuitai/internal/llms/openai"
	"github.com/Stosan/darksuitai/internal/memory/mongodb"
	"github.com/Stosan/darksuitai/internal/utilities"
)

var llm LLM

func (ai ConvAI) Chat(prompt string) (string, error) {
	kwargs := make([]map[string]interface{}, 5)
	for key := range ai.ModelType {
		switch key {
		case "openai":
			for k, v := range ai.ModelKwargs {
				kwargs[k] = map[string]interface{}{
					"model":          ai.ModelType["openai"],
					"max_tokens":     v.MaxTokens,
					"temperature":    v.Temperature,
					"stream":         v.Stream,
					"stop_sequences": v.StopSequences,
				}

			}
			llm = oai.ChatOAI(kwargs...)
		case "anthropic":
			for k, v := range ai.ModelKwargs {
				kwargs[k] = map[string]interface{}{
					"model":          ai.ModelType["anthropic"],
					"max_tokens":     v.MaxTokens,
					"temperature":    v.Temperature,
					"stream":         v.Stream,
					"stop_sequences": v.StopSequences,
				}

			}

			llm = ant.ChatAnth(kwargs...)
		default:
			llm = nil
		}
	}
	promptMap := ai.PromptKeys
	promptMap["query"] = []byte(prompt)
	promptMap["chat_history"] = []byte(mongodb.RetrieveMemoryWithK(ai.MongoDB, "", 6))
	promptTemplate := utilities.CustomFormat(ai.ChatInstruction, promptMap)
	resp, err := llm.StreamCompleteChat(string(promptTemplate), "")
	return resp, err
}

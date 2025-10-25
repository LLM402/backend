package service

import (
	"github.com/QuantumNous/new-api/dto"
)













func ResponseText2Usage(responseText string, modeName string, promptTokens int) *dto.Usage {
	usage := &dto.Usage{}
	usage.PromptTokens = promptTokens
	ctkm := CountTextToken(responseText, modeName)
	usage.CompletionTokens = ctkm
	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	return usage
}

func ValidUsage(usage *dto.Usage) bool {
	return usage != nil && (usage.PromptTokens != 0 || usage.CompletionTokens != 0)
}

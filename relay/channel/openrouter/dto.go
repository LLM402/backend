package openrouter

import "encoding/json"

type RequestReasoning struct {
	
	Effort    string `json:"effort,omitempty"`     
	MaxTokens int    `json:"max_tokens,omitempty"` 
	
	Exclude bool `json:"exclude,omitempty"` 
}

type OpenRouterEnterpriseResponse struct {
	Data    json.RawMessage `json:"data"`
	Success bool            `json:"success"`
}

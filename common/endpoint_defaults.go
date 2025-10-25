package common

import "github.com/QuantumNous/new-api/constant"









type EndpointInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}


var defaultEndpointInfoMap = map[constant.EndpointType]EndpointInfo{
	constant.EndpointTypeOpenAI:          {Path: "/v1/chat/completions", Method: "POST"},
	constant.EndpointTypeOpenAIResponse:  {Path: "/v1/responses", Method: "POST"},
	constant.EndpointTypeAnthropic:       {Path: "/v1/messages", Method: "POST"},
	constant.EndpointTypeGemini:          {Path: "/v1beta/models/{model}:generateContent", Method: "POST"},
	constant.EndpointTypeJinaRerank:      {Path: "/rerank", Method: "POST"},
	constant.EndpointTypeImageGeneration: {Path: "/v1/images/generations", Method: "POST"},
	constant.EndpointTypeEmbeddings:      {Path: "/v1/embeddings", Method: "POST"},
}


func GetDefaultEndpointInfo(et constant.EndpointType) (EndpointInfo, bool) {
	info, ok := defaultEndpointInfoMap[et]
	return info, ok
}

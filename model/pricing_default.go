package model

import (
	"strings"
)


var defaultVendorRules = map[string]string{
	"gpt":      "OpenAI",
	"dall-e":   "OpenAI",
	"whisper":  "OpenAI",
	"o1":       "OpenAI",
	"o3":       "OpenAI",
	"claude":   "Anthropic",
	"gemini":   "Google",
	"moonshot": "Moonshot",
	"kimi":     "Moonshot",
	"chatglm":  "Smart Spectrum",
	"glm-":     "Smart Spectrum",
	"qwen":     "Alibaba",
	"deepseek": "DeepSeek",
	"abab":     "MiniMax",
	"ernie":    "Baidu",
	"spark":    "iFlytek",
	"hunyuan":  "Tencent",
	"command":  "Cohere",
	"@cf/":     "Cloudflare",
	"360":      "360",
	"yi":       "Everything in the world",
	"jina":     "Jina",
	"mistral":  "Mistral",
	"grok":     "xAI",
	"llama":    "Meta",
	"doubao":   "ByteDance",
	"kling":    "Kuaishou",
	"jimeng":   "Dream",
	"vidu":     "Vidu",
}


var defaultVendorIcons = map[string]string{
	"OpenAI":     "OpenAI",
	"Anthropic":  "Claude.Color",
	"Google":     "Gemini.Color",
	"Moonshot":   "Moonshot",
	"Smart Spectrum":         "Zhipu.Color",
	"Alibaba":       "Qwen.Color",
	"DeepSeek":   "DeepSeek.Color",
	"MiniMax":    "Minimax.Color",
	"Baidu":         "Wenxin.Color",
	"iFlytek":         "Spark.Color",
	"Tencent":         "Hunyuan.Color",
	"Cohere":     "Cohere.Color",
	"Cloudflare": "Cloudflare.Color",
	"360":        "Ai360.Color",
	"Everything in the world":       "Yi.Color",
	"Jina":       "Jina",
	"Mistral":    "Mistral.Color",
	"xAI":        "XAI",
	"Meta":       "Ollama",
	"ByteDance":       "Doubao.Color",
	"Kuaishou":         "Kling.Color",
	"Dream":         "Jimeng.Color",
	"Vidu":       "Vidu",
	"Microsoft":         "AzureAI",
	"Microsoft":  "AzureAI",
	"Azure":      "AzureAI",
}


func initDefaultVendorMapping(metaMap map[string]*Model, vendorMap map[int]*Vendor, enableAbilities []AbilityWithChannel) {
	for _, ability := range enableAbilities {
		modelName := ability.Model
		if _, exists := metaMap[modelName]; exists {
			continue
		}

		
		vendorID := 0
		modelLower := strings.ToLower(modelName)
		for pattern, vendorName := range defaultVendorRules {
			if strings.Contains(modelLower, pattern) {
				vendorID = getOrCreateVendor(vendorName, vendorMap)
				break
			}
		}

		
		metaMap[modelName] = &Model{
			ModelName: modelName,
			VendorID:  vendorID,
			Status:    1,
			NameRule:  NameRuleExact,
		}
	}
}


func getOrCreateVendor(vendorName string, vendorMap map[int]*Vendor) int {
	
	for id, vendor := range vendorMap {
		if vendor.Name == vendorName {
			return id
		}
	}

	
	newVendor := &Vendor{
		Name:   vendorName,
		Status: 1,
		Icon:   getDefaultVendorIcon(vendorName),
	}

	if err := newVendor.Insert(); err != nil {
		return 0
	}

	vendorMap[newVendor.Id] = newVendor
	return newVendor.Id
}


func getDefaultVendorIcon(vendorName string) string {
	if icon, exists := defaultVendorIcons[vendorName]; exists {
		return icon
	}
	return ""
}

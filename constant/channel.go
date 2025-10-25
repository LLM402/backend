package constant

const (
	ChannelTypeUnknown        = 0
	ChannelTypeOpenAI         = 1
	ChannelTypeMidjourney     = 2
	ChannelTypeAzure          = 3
	ChannelTypeOllama         = 4
	ChannelTypeMidjourneyPlus = 5
	ChannelTypeOpenAIMax      = 6
	ChannelTypeOhMyGPT        = 7
	ChannelTypeCustom         = 8
	ChannelTypeAILS           = 9
	ChannelTypeAIProxy        = 10
	ChannelTypePaLM           = 11
	ChannelTypeAPI2GPT        = 12
	ChannelTypeAIGC2D         = 13
	ChannelTypeAnthropic      = 14
	ChannelTypeBaidu          = 15
	ChannelTypeZhipu          = 16
	ChannelTypeAli            = 17
	ChannelTypeXunfei         = 18
	ChannelType360            = 19
	ChannelTypeOpenRouter     = 20
	ChannelTypeAIProxyLibrary = 21
	ChannelTypeFastGPT        = 22
	ChannelTypeTencent        = 23
	ChannelTypeGemini         = 24
	ChannelTypeMoonshot       = 25
	ChannelTypeZhipu_v4       = 26
	ChannelTypePerplexity     = 27
	ChannelTypeLingYiWanWu    = 31
	ChannelTypeAws            = 33
	ChannelTypeCohere         = 34
	ChannelTypeMiniMax        = 35
	ChannelTypeSunoAPI        = 36
	ChannelTypeDify           = 37
	ChannelTypeJina           = 38
	ChannelCloudflare         = 39
	ChannelTypeSiliconFlow    = 40
	ChannelTypeVertexAi       = 41
	ChannelTypeMistral        = 42
	ChannelTypeDeepSeek       = 43
	ChannelTypeMokaAI         = 44
	ChannelTypeVolcEngine     = 45
	ChannelTypeBaiduV2        = 46
	ChannelTypeXinference     = 47
	ChannelTypeXai            = 48
	ChannelTypeCoze           = 49
	ChannelTypeKling          = 50
	ChannelTypeJimeng         = 51
	ChannelTypeVidu           = 52
	ChannelTypeSubmodel       = 53
	ChannelTypeDoubaoVideo    = 54
	ChannelTypeSora           = 55
	ChannelTypeDummy          

)

var ChannelBaseURLs = []string{
	"",                                    
	"https://api.openai.com",              
	"https://oa.api2d.net",                
	"",                                    
	"http://localhost:11434",              
	"https://api.openai-sb.com",           
	"https://api.openaimax.com",           
	"https://api.ohmygpt.com",             
	"",                                    
	"https://api.caipacity.com",           
	"https://api.aiproxy.io",              
	"",                                    
	"https://api.api2gpt.com",             
	"https://api.aigc2d.com",              
	"https://api.anthropic.com",           
	"https://aip.baidubce.com",            
	"https://open.bigmodel.cn",            
	"https://dashscope.aliyuncs.com",      
	"",                                    
	"https://api.360.cn",                  
	"https://openrouter.ai/api",           
	"https://api.aiproxy.io",              
	"https://fastgpt.run/api/openapi",     
	"https://hunyuan.tencentcloudapi.com", 
	"https://generativelanguage.googleapis.com", 
	"https://api.moonshot.cn",                   
	"https://open.bigmodel.cn",                  
	"https://api.perplexity.ai",                 
	"",                                          
	"",                                          
	"",                                          
	"https://api.lingyiwanwu.com",               
	"",                                          
	"",                                          
	"https://api.cohere.ai",                     
	"https://api.minimax.chat",                  
	"",                                          
	"https://api.dify.ai",                       
	"https://api.jina.ai",                       
	"https://api.cloudflare.com",                
	"https://api.siliconflow.cn",                
	"",                                          
	"https://api.mistral.ai",                    
	"https://api.deepseek.com",                  
	"https://api.moka.ai",                       
	"https://ark.cn-beijing.volces.com",         
	"https://qianfan.baidubce.com",              
	"",                                          
	"https://api.x.ai",                          
	"https://api.coze.cn",                       
	"https://api.klingai.com",                   
	"https://visual.volcengineapi.com",          
	"https://api.vidu.cn",                       
	"https://llm.submodel.ai",                   
	"https://ark.cn-beijing.volces.com",         
	"https://api.openai.com",                    
}

var ChannelTypeNames = map[int]string{
	ChannelTypeUnknown:        "Unknown",
	ChannelTypeOpenAI:         "OpenAI",
	ChannelTypeMidjourney:     "Midjourney",
	ChannelTypeAzure:          "Azure",
	ChannelTypeOllama:         "Ollama",
	ChannelTypeMidjourneyPlus: "MidjourneyPlus",
	ChannelTypeOpenAIMax:      "OpenAIMax",
	ChannelTypeOhMyGPT:        "OhMyGPT",
	ChannelTypeCustom:         "Custom",
	ChannelTypeAILS:           "AILS",
	ChannelTypeAIProxy:        "AIProxy",
	ChannelTypePaLM:           "PaLM",
	ChannelTypeAPI2GPT:        "API2GPT",
	ChannelTypeAIGC2D:         "AIGC2D",
	ChannelTypeAnthropic:      "Anthropic",
	ChannelTypeBaidu:          "Baidu",
	ChannelTypeZhipu:          "Zhipu",
	ChannelTypeAli:            "Ali",
	ChannelTypeXunfei:         "Xunfei",
	ChannelType360:            "360",
	ChannelTypeOpenRouter:     "OpenRouter",
	ChannelTypeAIProxyLibrary: "AIProxyLibrary",
	ChannelTypeFastGPT:        "FastGPT",
	ChannelTypeTencent:        "Tencent",
	ChannelTypeGemini:         "Gemini",
	ChannelTypeMoonshot:       "Moonshot",
	ChannelTypeZhipu_v4:       "ZhipuV4",
	ChannelTypePerplexity:     "Perplexity",
	ChannelTypeLingYiWanWu:    "LingYiWanWu",
	ChannelTypeAws:            "AWS",
	ChannelTypeCohere:         "Cohere",
	ChannelTypeMiniMax:        "MiniMax",
	ChannelTypeSunoAPI:        "SunoAPI",
	ChannelTypeDify:           "Dify",
	ChannelTypeJina:           "Jina",
	ChannelCloudflare:         "Cloudflare",
	ChannelTypeSiliconFlow:    "SiliconFlow",
	ChannelTypeVertexAi:       "VertexAI",
	ChannelTypeMistral:        "Mistral",
	ChannelTypeDeepSeek:       "DeepSeek",
	ChannelTypeMokaAI:         "MokaAI",
	ChannelTypeVolcEngine:     "VolcEngine",
	ChannelTypeBaiduV2:        "BaiduV2",
	ChannelTypeXinference:     "Xinference",
	ChannelTypeXai:            "xAI",
	ChannelTypeCoze:           "Coze",
	ChannelTypeKling:          "Kling",
	ChannelTypeJimeng:         "Jimeng",
	ChannelTypeVidu:           "Vidu",
	ChannelTypeSubmodel:       "Submodel",
	ChannelTypeDoubaoVideo:    "DoubaoVideo",
	ChannelTypeSora:           "Sora",
}

func GetChannelTypeName(channelType int) string {
	if name, ok := ChannelTypeNames[channelType]; ok {
		return name
	}
	return "Unknown"
}

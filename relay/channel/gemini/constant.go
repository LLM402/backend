package gemini

var ModelList = []string{
	
	"gemini-1.5-pro", "gemini-1.5-flash", "gemini-1.5-flash-8b",
	"gemini-2.0-flash",
	
	"gemini-1.5-pro-latest", "gemini-1.5-flash-latest",
	
	"gemini-2.0-flash-lite-preview",
	
	"gemini-exp-1206",
	
	"gemini-2.0-flash-exp",
	
	"gemini-2.0-pro-exp",
	
	"gemini-2.0-flash-thinking-exp",
	"gemini-2.5-pro-exp-03-25",
	"gemini-2.5-pro-preview-03-25",
	
	"imagen-3.0-generate-002",
	
	"gemini-embedding-exp-03-07",
	"text-embedding-004",
	"embedding-001",
}

var SafetySettingList = []string{
	"HARM_CATEGORY_HARASSMENT",
	"HARM_CATEGORY_HATE_SPEECH",
	"HARM_CATEGORY_SEXUALLY_EXPLICIT",
	"HARM_CATEGORY_DANGEROUS_CONTENT",
	"HARM_CATEGORY_CIVIC_INTEGRITY",
}

var ChannelName = "google gemini"

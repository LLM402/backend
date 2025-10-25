package types

type FileType string

const (
	FileTypeImage FileType = "image" 
	FileTypeAudio FileType = "audio" 
	FileTypeVideo FileType = "video" 
	FileTypeFile  FileType = "file"  
)

type TokenType string

const (
	TokenTypeTextNumber TokenType = "text_number" 
	TokenTypeTokenizer  TokenType = "tokenizer"   
	TokenTypeImage      TokenType = "image"       
)

type TokenCountMeta struct {
	TokenType     TokenType   `json:"token_type,omitempty"`     
	CombineText   string      `json:"combine_text,omitempty"`   
	ToolsCount    int         `json:"tools_count,omitempty"`    
	NameCount     int         `json:"name_count,omitempty"`     
	MessagesCount int         `json:"messages_count,omitempty"` 
	Files         []*FileMeta `json:"files,omitempty"`          
	MaxTokens     int         `json:"max_tokens,omitempty"`     

	ImagePriceRatio float64 `json:"image_ratio,omitempty"` 
	
}

type FileMeta struct {
	FileType
	MimeType   string
	OriginData string 
	Detail     string
	ParsedData *LocalFileData
}

type RequestMeta struct {
	OriginalModelName string `json:"original_model_name"`
	UserUsingGroup    string `json:"user_using_group"`
	PromptTokens      int    `json:"prompt_tokens"`
	PreConsumedQuota  int    `json:"pre_consumed_quota"`
}

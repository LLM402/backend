package dto

type VideoRequest struct {
	Model          string         `json:"model,omitempty" example:"kling-v1"`                                                                                                                                    
	Prompt         string         `json:"prompt,omitempty" example:"The astronaut stood up and walked."`                                                                                                                                   
	Image          string         `json:"image,omitempty" example:"https://h2.inkwai.com/bs2/upload-ylab-stunt/se/ai_portal_queue_mmu_image_upscale_aiweb/3214b798-e1b4-4b00-b7af-72b5b0417420_raw_image_0.jpg"` 
	Duration       float64        `json:"duration" example:"5.0"`                                                                                                                                                
	Width          int            `json:"width" example:"512"`                                                                                                                                                   
	Height         int            `json:"height" example:"512"`                                                                                                                                                  
	Fps            int            `json:"fps,omitempty" example:"30"`                                                                                                                                            
	Seed           int            `json:"seed,omitempty" example:"20231234"`                                                                                                                                     
	N              int            `json:"n,omitempty" example:"1"`                                                                                                                                               
	ResponseFormat string         `json:"response_format,omitempty" example:"url"`                                                                                                                               
	User           string         `json:"user,omitempty" example:"user-1234"`                                                                                                                                    
	Metadata       map[string]any `json:"metadata,omitempty"`                                                                                                                                                    
}


type VideoResponse struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
}


type VideoTaskResponse struct {
	TaskId   string             `json:"task_id" example:"abcd1234efgh"` 
	Status   string             `json:"status" example:"succeeded"`     
	Url      string             `json:"url,omitempty"`                  
	Format   string             `json:"format,omitempty" example:"mp4"` 
	Metadata *VideoTaskMetadata `json:"metadata,omitempty"`             
	Error    *VideoTaskError    `json:"error,omitempty"`                
}


type VideoTaskMetadata struct {
	Duration float64 `json:"duration" example:"5.0"`  
	Fps      int     `json:"fps" example:"30"`        
	Width    int     `json:"width" example:"512"`     
	Height   int     `json:"height" example:"512"`    
	Seed     int     `json:"seed" example:"20231234"` 
}


type VideoTaskError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

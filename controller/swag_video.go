package controller

import (
	"github.com/gin-gonic/gin"
)

















func VideoGenerations(c *gin.Context) {
}















func VideoGenerationsTaskId(c *gin.Context) {
}















func KlingText2VideoGenerations(c *gin.Context) {
}

type KlingText2VideoRequest struct {
	ModelName      string              `json:"model_name,omitempty" example:"kling-v1"`
	Prompt         string              `json:"prompt" binding:"required" example:"A cat playing piano in the garden"`
	NegativePrompt string              `json:"negative_prompt,omitempty" example:"blurry, low quality"`
	CfgScale       float64             `json:"cfg_scale,omitempty" example:"0.7"`
	Mode           string              `json:"mode,omitempty" example:"std"`
	CameraControl  *KlingCameraControl `json:"camera_control,omitempty"`
	AspectRatio    string              `json:"aspect_ratio,omitempty" example:"16:9"`
	Duration       string              `json:"duration,omitempty" example:"5"`
	CallbackURL    string              `json:"callback_url,omitempty" example:"https://your.domain/callback"`
	ExternalTaskId string              `json:"external_task_id,omitempty" example:"custom-task-001"`
}

type KlingCameraControl struct {
	Type   string             `json:"type,omitempty" example:"simple"`
	Config *KlingCameraConfig `json:"config,omitempty"`
}

type KlingCameraConfig struct {
	Horizontal float64 `json:"horizontal,omitempty" example:"2.5"`
	Vertical   float64 `json:"vertical,omitempty" example:"0"`
	Pan        float64 `json:"pan,omitempty" example:"0"`
	Tilt       float64 `json:"tilt,omitempty" example:"0"`
	Roll       float64 `json:"roll,omitempty" example:"0"`
	Zoom       float64 `json:"zoom,omitempty" example:"0"`
}















func KlingImage2VideoGenerations(c *gin.Context) {
}

type KlingImage2VideoRequest struct {
	ModelName      string              `json:"model_name,omitempty" example:"kling-v2-master"`
	Image          string              `json:"image" binding:"required" example:"https://h2.inkwai.com/bs2/upload-ylab-stunt/se/ai_portal_queue_mmu_image_upscale_aiweb/3214b798-e1b4-4b00-b7af-72b5b0417420_raw_image_0.jpg"`
	Prompt         string              `json:"prompt,omitempty" example:"A cat playing piano in the garden"`
	NegativePrompt string              `json:"negative_prompt,omitempty" example:"blurry, low quality"`
	CfgScale       float64             `json:"cfg_scale,omitempty" example:"0.7"`
	Mode           string              `json:"mode,omitempty" example:"std"`
	CameraControl  *KlingCameraControl `json:"camera_control,omitempty"`
	AspectRatio    string              `json:"aspect_ratio,omitempty" example:"16:9"`
	Duration       string              `json:"duration,omitempty" example:"5"`
	CallbackURL    string              `json:"callback_url,omitempty" example:"https://your.domain/callback"`
	ExternalTaskId string              `json:"external_task_id,omitempty" example:"custom-task-002"`
}









func KlingImage2videoTaskId(c *gin.Context) {}









func KlingText2videoTaskId(c *gin.Context) {}

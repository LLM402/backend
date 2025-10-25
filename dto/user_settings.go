package dto

type UserSetting struct {
	NotifyType            string  `json:"notify_type,omitempty"`                    
	QuotaWarningThreshold float64 `json:"quota_warning_threshold,omitempty"`        
	WebhookUrl            string  `json:"webhook_url,omitempty"`                    
	WebhookSecret         string  `json:"webhook_secret,omitempty"`                 
	NotificationEmail     string  `json:"notification_email,omitempty"`             
	BarkUrl               string  `json:"bark_url,omitempty"`                       
	GotifyUrl             string  `json:"gotify_url,omitempty"`                     
	GotifyToken           string  `json:"gotify_token,omitempty"`                   
	GotifyPriority        int     `json:"gotify_priority"`                          
	AcceptUnsetRatioModel bool    `json:"accept_unset_model_ratio_model,omitempty"` 
	RecordIpLog           bool    `json:"record_ip_log,omitempty"`                  
	SidebarModules        string  `json:"sidebar_modules,omitempty"`                
}

var (
	NotifyTypeEmail   = "email"   
	NotifyTypeWebhook = "webhook" 
	NotifyTypeBark    = "bark"    
	NotifyTypeGotify  = "gotify"  
)

package common

import (
	
	
	"sync"
	"time"

	"github.com/google/uuid"
)

var StartTime = time.Now().Unix() 
var Version = "v0.0.0"            
var SystemName = "New API"
var Footer = ""
var Logo = ""
var TopUpLink = ""



var QuotaPerUnit = 500 * 1000.0 

var DisplayInCurrencyEnabled = true
var DisplayTokenStatEnabled = true
var DrawingEnabled = true
var TaskEnabled = true
var DataExportEnabled = true
var DataExportInterval = 5         
var DataExportDefaultTime = "hour" 
var DefaultCollapseSidebar = false 


var EmailAliasRestrictionEnabled = false  
var BatchUpdateEnabled = false
var RelayTimeout int 
var GeminiSafetySetting string


var CohereSafetySetting string's unit is seconds
// Shouldn't larger then RateLimitKeyExpirationDuration
var (
	GlobalApiRateLimitEnable   bool
	GlobalApiRateLimitNum      int
	GlobalApiRateLimitDuration int64

	GlobalWebRateLimitEnable   bool
	GlobalWebRateLimitNum      int
	GlobalWebRateLimitDuration int64

	UploadRateLimitNum            = 10
	UploadRateLimitDuration int64 = 60

	DownloadRateLimitNum            = 10
	DownloadRateLimitDuration int64 = 60

	CriticalRateLimitNum            = 20
	CriticalRateLimitDuration int64 = 20 * 60
)

var RateLimitKeyExpirationDuration = 20 * time.Minute

const (
	UserStatusEnabled  = 1 
	UserStatusDisabled = 2 
)

const (
	TokenStatusEnabled   = 1 
	TokenStatusDisabled  = 2 
	TokenStatusExpired   = 3
	TokenStatusExhausted = 4
)

const (
	RedemptionCodeStatusEnabled  = 1 
	RedemptionCodeStatusDisabled = 2 
	RedemptionCodeStatusUsed     = 3 
)

const (
	ChannelStatusUnknown          = 0
	ChannelStatusEnabled          = 1 
	ChannelStatusManuallyDisabled = 2 
	ChannelStatusAutoDisabled     = 3
)

const (
	TopUpStatusPending = "pending"
	TopUpStatusSuccess = "success"
	TopUpStatusExpired = "expired"
)

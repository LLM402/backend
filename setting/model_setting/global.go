package model_setting

import (
	"github.com/QuantumNous/new-api/setting/config"
)

type GlobalSettings struct {
	PassThroughRequestEnabled bool `json:"pass_through_request_enabled"`
}


var defaultOpenaiSettings = GlobalSettings{
	PassThroughRequestEnabled: false,
}


var globalSettings = defaultOpenaiSettings

func init() {
	
	config.GlobalConfig.Register("global", &globalSettings)
}

func GetGlobalSettings() *GlobalSettings {
	return &globalSettings
}

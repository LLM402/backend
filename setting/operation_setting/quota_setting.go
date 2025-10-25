package operation_setting

import "github.com/QuantumNous/new-api/setting/config"

type QuotaSetting struct {
	EnableFreeModelPreConsume bool `json:"enable_free_model_pre_consume"` 
}


var quotaSetting = QuotaSetting{
	EnableFreeModelPreConsume: true,
}

func init() {
	
	config.GlobalConfig.Register("quota_setting", &quotaSetting)
}

func GetQuotaSetting() *QuotaSetting {
	return &quotaSetting
}

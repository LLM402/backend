

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)


func MigrateConsoleSetting(c *gin.Context) {
	
	opts, err := model.AllOption()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	
	valMap := map[string]string{}
	for _, o := range opts {
		valMap[o.Key] = o.Value
	}

	
	if v := valMap["ApiInfo"]; v != "" {
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			if len(arr) > 50 {
				arr = arr[:50]
			}
			bytes, _ := json.Marshal(arr)
			model.UpdateOption("console_setting.api_info", string(bytes))
		}
		model.UpdateOption("ApiInfo", "")
	}
	
	if v := valMap["Announcements"]; v != "" {
		model.UpdateOption("console_setting.announcements", v)
		model.UpdateOption("Announcements", "")
	}
	
	if v := valMap["FAQ"]; v != "" {
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			out := []map[string]interface{}{}
			for _, item := range arr {
				q, _ := item["question"].(string)
				if q == "" {
					q, _ = item["title"].(string)
				}
				a, _ := item["answer"].(string)
				if a == "" {
					a, _ = item["content"].(string)
				}
				if q != "" && a != "" {
					out = append(out, map[string]interface{}{"question": q, "answer": a})
				}
			}
			if len(out) > 50 {
				out = out[:50]
			}
			bytes, _ := json.Marshal(out)
			model.UpdateOption("console_setting.faq", string(bytes))
		}
		model.UpdateOption("FAQ", "")
	}
	
	url := valMap["UptimeKumaUrl"]
	slug := valMap["UptimeKumaSlug"]
	if url != "" && slug != "" {
		
		groups := []map[string]interface{}{
			{
				"id":           1,
				"categoryName": "old",
				"url":          url,
				"slug":         slug,
				"description":  "",
			},
		}
		bytes, _ := json.Marshal(groups)
		model.UpdateOption("console_setting.uptime_kuma_groups", string(bytes))
	}
	
	if url != "" {
		model.UpdateOption("UptimeKumaUrl", "")
	}
	if slug != "" {
		model.UpdateOption("UptimeKumaSlug", "")
	}

	
	oldKeys := []string{"ApiInfo", "Announcements", "FAQ", "UptimeKumaUrl", "UptimeKumaSlug"}
	model.DB.Where("key IN ?", oldKeys).Delete(&model.Option{})

	
	model.InitOptionMap()
	common.SysLog("console setting migrated")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "migrated"})
}

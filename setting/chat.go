package setting

import (
	"encoding/json"

	"github.com/QuantumNous/new-api/common"
)

var Chats = []map[string]string{
	
	
	
	{
		"Cherry Studio": "cherrystudio://providers/api-keys?v=1&data={cherryConfig}",
	},
	{
		"Smooth reading": "fluentread",
	},
	{
		"Lobe Chat Official Example": "https:
	},
	{
		"AI as Workspace": "https:
	},
	{
		"AMA asks the sky": "ama://set-api-key?server={address}&key={key}",
	},
	{
		"OpenCat": "opencat://team/join?domain={address}&token={key}",
	},
}

func UpdateChatsByJsonString(jsonString string) error {
	Chats = make([]map[string]string, 0)
	return json.Unmarshal([]byte(jsonString), &Chats)
}

func Chats2JsonString() string {
	jsonBytes, err := json.Marshal(Chats)
	if err != nil {
		common.SysLog("error marshalling chats: " + err.Error())
		return "[]"
	}
	return string(jsonBytes)
}

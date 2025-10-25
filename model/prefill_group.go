package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/QuantumNous/new-api/common"

	"gorm.io/gorm"
)









type JSONValue json.RawMessage


func (j JSONValue) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}


func (j *JSONValue) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*j = nil
		return nil
	case []byte:
		
		b := make([]byte, len(v))
		copy(b, v)
		*j = JSONValue(b)
		return nil
	case string:
		*j = JSONValue([]byte(v))
		return nil
	default:
		
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		*j = JSONValue(b)
		return nil
	}
}


func (j JSONValue) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}


func (j *JSONValue) UnmarshalJSON(data []byte) error {
	if data == nil {
		*j = nil
		return nil
	}
	b := make([]byte, len(data))
	copy(b, data)
	*j = JSONValue(b)
	return nil
}

type PrefillGroup struct {
	Id          int            `json:"id"`
	Name        string         `json:"name" gorm:"size:64;not null;uniqueIndex:uk_prefill_name,where:deleted_at IS NULL"`
	Type        string         `json:"type" gorm:"size:32;index;not null"`
	Items       JSONValue      `json:"items" gorm:"type:json"`
	Description string         `json:"description,omitempty" gorm:"type:varchar(255)"`
	CreatedTime int64          `json:"created_time" gorm:"bigint"`
	UpdatedTime int64          `json:"updated_time" gorm:"bigint"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}


func (g *PrefillGroup) Insert() error {
	now := common.GetTimestamp()
	g.CreatedTime = now
	g.UpdatedTime = now
	return DB.Create(g).Error
}


func IsPrefillGroupNameDuplicated(id int, name string) (bool, error) {
	if name == "" {
		return false, nil
	}
	var cnt int64
	err := DB.Model(&PrefillGroup{}).Where("name = ? AND id <> ?", name, id).Count(&cnt).Error
	return cnt > 0, err
}


func (g *PrefillGroup) Update() error {
	g.UpdatedTime = common.GetTimestamp()
	return DB.Save(g).Error
}


func DeletePrefillGroupByID(id int) error {
	return DB.Delete(&PrefillGroup{}, id).Error
}


func GetAllPrefillGroups(groupType string) ([]*PrefillGroup, error) {
	var groups []*PrefillGroup
	query := DB.Model(&PrefillGroup{})
	if groupType != "" {
		query = query.Where("type = ?", groupType)
	}
	if err := query.Order("updated_time DESC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

package widget

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/serializable_meta"
)

type QorWidgetSetting struct {
	gorm.Model
	Scope string
	Key   string
	serializable_meta.SerializableMeta
}

func findSettingByNameAndKey(db *gorm.DB, kind string, key string) *QorWidgetSetting {
	setting := QorWidgetSetting{}
	if db.Where("`key` = ? AND kind = ?", key, kind).First(&setting).RecordNotFound() {
		setting.Key = key
		setting.Kind = kind
		db.Save(&setting)
	}
	return &setting
}

// GetSerializableArgumentResource get setting's argument's resource
func (setting *QorWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	widget, _ := GetWidget(setting.Kind)
	return widget.Setting
}

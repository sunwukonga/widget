package widget

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/serializable_meta"
)

// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	gorm.Model
	Template string
	Scope    string
	Name     string
	serializable_meta.SerializableMeta
}

func findSettingByNameAndKinds(db *gorm.DB, widgetKey string, widgetName string) *QorWidgetSetting {
	setting := QorWidgetSetting{}
	if db.Where("name = ? AND kind = ?", widgetKey, widgetName).First(&setting).RecordNotFound() {
		setting.Name = widgetKey
		setting.Kind = widgetName
		db.Save(&setting)
	}
	return &setting
}

// GetSerializableArgumentResource get setting's argument's resource
func (setting *QorWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	return GetWidget(setting.Kind).Setting
}

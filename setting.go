package widget

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/serializable_meta"
)

// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	Name     string `gorm:"primary_key"`
	Scope    string `gorm:"primary_key;default:'default'"`
	Template string
	serializable_meta.SerializableMeta
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetTemplate get used widget template
func (qorWidgetSetting QorWidgetSetting) GetTemplate() string {
	if widget := GetWidget(qorWidgetSetting.Kind); widget != nil {
		for _, value := range widget.Templates {
			if value == qorWidgetSetting.Template {
				return value
			}
		}

		// return first value of defined widget templates
		for _, value := range widget.Templates {
			return value
		}
	}
	return ""
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

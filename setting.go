package widget

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/serializable_meta"
)

// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	Name      string `gorm:"primary_key"`
	Scope     string `gorm:"primary_key;default:'default'"`
	GroupName string
	Template  string
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

func findSettingByName(db *gorm.DB, widgetName string, scopes []string, widgetsGroupNameOrWidgetName string) *QorWidgetSetting {
	var setting *QorWidgetSetting
	var settings []QorWidgetSetting

	db.Where("name = ? AND scope IN (?)", widgetName, append(scopes, "default")).Find(&settings)

	if len(settings) > 0 {
	OUTTER:
		for _, scope := range scopes {
			for _, s := range settings {
				if s.Scope == scope {
					setting = &s
					break OUTTER
				}
			}
		}
	}

	// use default setting
	if setting == nil {
		for _, s := range settings {
			if s.Scope == "default" {
				setting = &s
			}
		}
	}

	if setting == nil {
		setting = &QorWidgetSetting{Name: widgetName, Scope: "default"}
		setting.Kind = widgetsGroupNameOrWidgetName
		db.Create(setting)
	}

	return setting
}

// GetSerializableArgumentResource get setting's argument's resource
func (qorWidgetSetting *QorWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	return GetWidget(qorWidgetSetting.Kind).Setting
}

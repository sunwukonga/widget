package widget

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/roles"
	"github.com/qor/serializable_meta"
)

// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	Name        string `gorm:"primary_key"`
	WidgetType  string `gorm:"primary_key"`
	Scope       string `gorm:"primary_key;default:'default'"`
	GroupName   string
	ActivatedAt *time.Time
	Template    string
	serializable_meta.SerializableMeta
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (widgetSetting *QorWidgetSetting) ResourceName() string {
	return "Widget Setting"
}

func (widgetSetting *QorWidgetSetting) BeforeCreate() {
	now := time.Now()
	widgetSetting.ActivatedAt = &now
}

func (widgetSetting *QorWidgetSetting) GetSerializableArgumentKind() string {
	if widgetSetting.WidgetType != "" {
		return widgetSetting.WidgetType
	}
	return widgetSetting.Kind
}

func (widgetSetting *QorWidgetSetting) SetSerializableArgumentKind(name string) {
	widgetSetting.WidgetType = name
	widgetSetting.Kind = name
}

// GetTemplate get used widget template
func (qorWidgetSetting QorWidgetSetting) GetTemplate() string {
	if widget := GetWidget(qorWidgetSetting.GetSerializableArgumentKind()); widget != nil {
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

	db.Where("name = ? AND scope IN (?)", widgetName, append(scopes, "default")).Order("activated_at DESC").Find(&settings)

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
				break
			}
		}
	}

	if setting == nil {
		setting = &QorWidgetSetting{Name: widgetName, Scope: "default"}
		setting.GroupName = widgetsGroupNameOrWidgetName
		setting.SetSerializableArgumentKind(widgetsGroupNameOrWidgetName)
		db.Create(setting)
	} else if setting.GroupName != widgetsGroupNameOrWidgetName {
		setting.GroupName = widgetsGroupNameOrWidgetName
		db.Save(setting)
	}

	return setting
}

// GetSerializableArgumentResource get setting's argument's resource
func (qorWidgetSetting *QorWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	return GetWidget(qorWidgetSetting.GetSerializableArgumentKind()).Setting
}

// ConfigureQorResource a method used to config Widget for qor admin
func (qorWidgetSetting *QorWidgetSetting) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.Meta(&admin.Meta{Name: "Name", Permission: roles.Deny(roles.Update, roles.Anyone)})

		res.Meta(&admin.Meta{
			Name: "ActivatedAt",
			Type: "hidden",
			Valuer: func(result interface{}, context *qor.Context) interface{} {
				return time.Now()
			},
		})

		res.Meta(&admin.Meta{
			Name: "Scope",
			Type: "hidden",
			Valuer: func(result interface{}, context *qor.Context) interface{} {
				if scope := context.Request.URL.Query().Get("widget_scope"); scope != "" {
					return scope
				}

				if setting, ok := result.(*QorWidgetSetting); ok {
					if setting.Scope != "" {
						return setting.Scope
					}
				}

				return "default"
			},
		})

		res.Meta(&admin.Meta{
			Name: "Widgets",
			Type: "select_one",
			Valuer: func(result interface{}, context *qor.Context) interface{} {
				if typ := context.Request.URL.Query().Get("widget_type"); typ != "" {
					return typ
				}

				if setting, ok := result.(*QorWidgetSetting); ok {
					return GetWidget(setting.GetSerializableArgumentKind()).Name
				}

				return ""
			},
			Collection: func(result interface{}, context *qor.Context) (results [][]string) {
				if setting, ok := result.(*QorWidgetSetting); ok {
					for _, group := range registeredWidgetsGroup {
						if group.Name == setting.GroupName {
							for _, widget := range group.Widgets {
								results = append(results, []string{widget, widget})
							}
						}
					}

					if len(results) == 0 {
						results = append(results, []string{setting.GetSerializableArgumentKind(), setting.GetSerializableArgumentKind()})
					}
				}
				return
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *qor.Context) {
				if setting, ok := result.(*QorWidgetSetting); ok {
					setting.SetSerializableArgumentKind(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name: "Template",
			Type: "select_one",
			Valuer: func(result interface{}, context *qor.Context) interface{} {
				if setting, ok := result.(*QorWidgetSetting); ok {
					return setting.GetTemplate()
				}
				return ""
			},
			Collection: func(result interface{}, context *qor.Context) (results [][]string) {
				if setting, ok := result.(*QorWidgetSetting); ok {
					if widget := GetWidget(setting.GetSerializableArgumentKind()); widget != nil {
						for _, value := range widget.Templates {
							results = append(results, []string{value, value})
						}
					}
				}
				return
			},
		})

		res.UseTheme("widget")

		res.IndexAttrs("Name", "CreatedAt", "UpdatedAt")
		res.ShowAttrs("Name", "Scope", "WidgetType", "Template", "Value", "CreatedAt", "UpdatedAt")
		res.EditAttrs(
			"Scope", "ActivatedAt", "Widgets", "Template",
			&admin.Section{
				Title: "Settings",
				Rows:  [][]string{{"Kind"}, {"SerializableMeta"}},
			},
		)
	}
}

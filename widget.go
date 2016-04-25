package widget

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

var (
	root, _           = os.Getwd()
	viewPaths         []string
	registeredWidgets []*Widget
)

// Config widget config
type Config struct {
	DB    *gorm.DB
	Admin *admin.Admin
}

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
	registerViewPath(path.Join(root, "app/views/widgets"))
}

// New new widgets container
func New(config *Config) *Widgets {
	widgets := &Widgets{Config: config}
	widgets.RegisterViewPath("app/views/widgets")
	return widgets
}

// Widgets widgets container
type Widgets struct {
	Config                *Config
	Resource              *admin.Resource
	WidgetSettingResource *admin.Resource
}

// RegisterWidget register a new widget
func (widgets *Widgets) RegisterWidget(w *Widget) {
	registeredWidgets = append(registeredWidgets, w)
}

// ConfigureQorResource a method used to config Widget for qor admin
func (widgets *Widgets) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		// register view paths
		admin.RegisterViewPath("github.com/qor/widget/views")

		// set resources
		res.Name = "Widget"
		widgets.Resource = res

		// set setting resource
		if widgets.WidgetSettingResource == nil {
			widgets.WidgetSettingResource = res.GetAdmin().NewResource(&QorWidgetSetting{}, &admin.Config{Name: res.Name, Permission: roles.Deny(roles.Create, roles.Anyone)})

			widgets.WidgetSettingResource.Meta(&admin.Meta{Name: "Name", Permission: roles.Deny(roles.Update, roles.Anyone)})
			widgets.WidgetSettingResource.Meta(&admin.Meta{
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
			widgets.WidgetSettingResource.Meta(&admin.Meta{
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
						if widget := GetWidget(setting.Kind); widget != nil {
							for _, value := range widget.Templates {
								results = append(results, []string{value, value})
							}
						}
					}
					return
				}})

			widgets.WidgetSettingResource.IndexAttrs("ID", "Name", "Template", "Kind", "CreatedAt", "UpdatedAt")
			widgets.WidgetSettingResource.EditAttrs("ID", "Scope", "Template", &admin.Section{Title: "Settings", Rows: [][]string{[]string{"Kind"}, []string{"SerializableMeta"}}})
		}

		// use widget theme
		res.UseTheme("widget")
		widgets.WidgetSettingResource.UseTheme("widget")

		for funcName, fc := range funcMap {
			res.GetAdmin().RegisterFuncMap(funcName, fc)
		}

		// configure routes
		controller := widgetController{Widgets: widgets}
		router := res.GetAdmin().GetRouter()
		router.Get(res.ToParam(), controller.Index)
		router.Get(fmt.Sprintf("%v/inline-edit", res.ToParam()), controller.InlineEdit)
		router.Get(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Get(fmt.Sprintf("%v/%v/edit", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Put(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Update)
	}
}

// Widget widget struct
type Widget struct {
	Name      string
	Templates []string
	Setting   *admin.Resource
	Context   func(context *Context, setting interface{}) *Context
}

// GetWidget get widget by name
func GetWidget(name string) *Widget {
	for _, w := range registeredWidgets {
		if w.Name == name {
			return w
		}
	}
	return nil
}

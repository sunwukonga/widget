package widget

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
)

var (
	root, _                = os.Getwd()
	viewPaths              []string
	registeredWidgets      []*Widget
	registeredWidgetsGroup []*WidgetsGroup
)

// Config widget config
type Config struct {
	DB            *gorm.DB
	Admin         *admin.Admin
	PreviewAssets []string
}

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
	registerViewPath(path.Join(root, "app/views/widgets"))
}

// New new widgets container
func New(config *Config) *Widgets {
	widgets := &Widgets{Config: config, funcMaps: template.FuncMap{}}
	widgets.RegisterViewPath("app/views/widgets")
	return widgets
}

// Widgets widgets container
type Widgets struct {
	funcMaps              template.FuncMap
	Config                *Config
	Resource              *admin.Resource
	WidgetSettingResource *admin.Resource
}

// RegisterWidget register a new widget
func (widgets *Widgets) RegisterWidget(w *Widget) {
	registeredWidgets = append(registeredWidgets, w)
}

// RegisterWidgetGroup register widgets group
func (widgets *Widgets) RegisterWidgetsGroup(group *WidgetsGroup) {
	registeredWidgetsGroup = append(registeredWidgetsGroup, group)
}

// // RegisterFuncMap register view funcs, it could be used when render templates
func (widgets *Widgets) RegisterFuncMap(name string, fc interface{}) {
	widgets.funcMaps[name] = fc
}

// ConfigureQorResource a method used to config Widget for qor admin
func (widgets *Widgets) ConfigureQorResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		// register view paths
		res.GetAdmin().RegisterViewPath("github.com/qor/widget/views")

		// set resources
		widgets.Resource = res
		res.Config.Invisible = true

		// set setting resource
		if widgets.WidgetSettingResource == nil {
			widgets.WidgetSettingResource = res.GetAdmin().AddResource(&QorWidgetSetting{}, &admin.Config{Name: res.Name})
		}

		for funcName, fc := range funcMap {
			res.GetAdmin().RegisterFuncMap(funcName, fc)
		}

		// configure routes
		controller := widgetController{Widgets: widgets}
		router := res.GetAdmin().GetRouter()
		router.Get(widgets.WidgetSettingResource.ToParam(), controller.Index)
		router.Get(fmt.Sprintf("%v/new", widgets.WidgetSettingResource.ToParam()), controller.New)
		router.Get(fmt.Sprintf("%v/!setting", widgets.WidgetSettingResource.ToParam()), controller.Setting)
		router.Get(fmt.Sprintf("%v/%v", widgets.WidgetSettingResource.ToParam(), widgets.WidgetSettingResource.ParamIDName()), controller.Edit)
		router.Get(fmt.Sprintf("%v/%v/!preview", widgets.WidgetSettingResource.ToParam(), widgets.WidgetSettingResource.ParamIDName()), controller.Preview)
		router.Get(fmt.Sprintf("%v/%v/edit", widgets.WidgetSettingResource.ToParam(), widgets.WidgetSettingResource.ParamIDName()), controller.Edit)
		router.Put(fmt.Sprintf("%v/%v", widgets.WidgetSettingResource.ToParam(), widgets.WidgetSettingResource.ParamIDName()), controller.Update)
		router.Get(fmt.Sprintf("%v/inline-edit", res.ToParam()), controller.InlineEdit)
	}
}

// Widget widget struct
type Widget struct {
	Name      string
	Templates []string
	Setting   *admin.Resource
	Context   func(context *Context, setting interface{}) *Context
}

// WidgetsGroup widgets Group
type WidgetsGroup struct {
	Name    string
	Widgets []string
}

// GetWidget get widget by name
func GetWidget(name string) *Widget {
	for _, w := range registeredWidgets {
		if w.Name == name {
			return w
		}
	}

	for _, g := range registeredWidgetsGroup {
		if g.Name == name {
			for _, widgetName := range g.Widgets {
				return GetWidget(widgetName)
			}
		}
	}
	return nil
}

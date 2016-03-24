package widget

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"html/template"
	"os"
	"path"
	"strings"
)

var (
	root, _ = os.Getwd()
)
var registeredWidgets []*Widget
var viewPaths = []string{}

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

// ConfigureQorResource a method used to config Widget for qor admin
func (w *WidgetInstance) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		admin.RegisterViewPath("github.com/qor/widget/views")

		// configure routes
		router := res.GetAdmin().GetRouter()
		w.SettingResource = res.GetAdmin().NewResource(&QorWidgetSetting{})
		w.SettingResource.IndexAttrs("ID", "Kind", "Key")
		w.SettingResource.Name = res.Name
		controller := widgetController{WidgetInstance: w}
		router.Get(res.ToParam(), controller.Index)
		router.Get(fmt.Sprintf("%v/frontend-edit", res.ToParam()), controller.FronendEdit)
		router.Get(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Get(fmt.Sprintf("%v/%v/edit", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Put(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Update)
	}
}

type WidgetInstance struct {
	Config          *Config
	SettingResource *admin.Resource
}

func New(config *Config) *WidgetInstance {
	instance := &WidgetInstance{Config: config}
	instance.RegisterViewPath("app/views/widgets")
	return instance
}

func (widgetInstance *WidgetInstance) RegisterWidget(w *Widget) {
	registeredWidgets = append(registeredWidgets, w)
}

func (widgetInstance *WidgetInstance) IncludeAssetTag() template.HTML {
	return "<script src=\"/admin/assets/javascripts/widget.js?theme=widget\"></script>"
}

type Widget struct {
	Name     string
	Template string
	Setting  *admin.Resource
	Context  func(context *Context, setting interface{}) *Context
}

func (w *Widget) nameForClass() string {
	return strings.ToLower(strings.Replace(w.Name, " ", "-", -1))
}

func GetWidget(name string) (w Widget, err error) {
	for _, w := range registeredWidgets {
		if w.Name == name {
			return *w, nil
		}
	}
	return Widget{}, fmt.Errorf("Widget: failed to find widget %v", name)
}

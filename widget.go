package widget

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"os"
	"path"
)

var (
	root, _ = os.Getwd()
)

type Config struct {
	DB    *gorm.DB
	Admin *admin.Admin
}

type WidgetInstance struct {
	Config          *Config
	SettingResource *admin.Resource
}

// ConfigureQorResource a method used to config Widget for qor admin
func (w *WidgetInstance) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		// configure routes
		router := res.GetAdmin().GetRouter()
		w.SettingResource = res.GetAdmin().NewResource(&QorWidgetSetting{})
		w.SettingResource.IndexAttrs("ID", "Kind", "Key")
		w.SettingResource.Name = res.Name
		controller := widgetController{WidgetInstance: w}
		router.Get(res.ToParam(), controller.Index)
		router.Get(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Get(fmt.Sprintf("%v/%v/edit", res.ToParam(), res.ParamIDName()), controller.Edit)
		router.Put(fmt.Sprintf("%v/%v", res.ToParam(), res.ParamIDName()), controller.Update)
	}
}

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
	registerViewPath(path.Join(root, "app/views/widgets"))
}

func New(config *Config) *WidgetInstance {
	instance := &WidgetInstance{Config: config}
	instance.RegisterViewPath("app/views/widgets")
	return instance
}

type Widget struct {
	Name     string
	Template string
	Setting  *admin.Resource
	Context  func(context *Context, setting interface{}) *Context
}

var registeredWidgets []*Widget
var viewPaths = []string{}

func (widgetInstance *WidgetInstance) RegisterWidget(w *Widget) {
	registeredWidgets = append(registeredWidgets, w)
}

func GetWidget(name string) (w Widget, err error) {
	for _, w := range registeredWidgets {
		if w.Name == name {
			return *w, nil
		}
	}
	return Widget{}, fmt.Errorf("Widget: failed to find widget %v", name)
}

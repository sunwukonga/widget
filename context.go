package widget

import (
	"fmt"
	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/utils"
)

// NewContext new widget context
func NewContext(options map[string]interface{}) *Context {
	return &Context{
		Options: options,
	}
}

// Context widget context
type Context struct {
	Widgets          *Widgets
	DB               *gorm.DB
	InlineEdit       bool
	AvailableWidgets []string
	Options          map[string]interface{}
}

// Get get option with name
func (context Context) Get(name string) (interface{}, bool) {
	if value, ok := context.Options[name]; ok {
		return value, true
	}

	return nil, false
}

// Set set option by name
func (context *Context) Set(name string, value interface{}) {
	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}
	context.Options[name] = value
}

// GetDB set option by name
func (context *Context) GetDB() *gorm.DB {
	if context.DB != nil {
		return context.DB
	}
	return context.Widgets.Config.DB
}

func (context *Context) Render(widgetName string, widgetGroupName string) template.HTML {
	var (
		visibleScopes []string
		widgets       = context.Widgets
		db            = context.GetDB()
	)

	for _, scope := range registeredScopes {
		if scope.Visible(context) {
			visibleScopes = append(visibleScopes, scope.ToParam())
		}
	}

	if setting := findSettingByName(db, widgetName, visibleScopes, widgetGroupName); setting != nil {
		var (
			widgetObj     = GetWidget(setting.GetSerializableArgumentKind())
			widgetSetting = widgetObj.Context(context, setting.GetSerializableArgument(setting))
			url           = widgets.settingEditURL(setting)
		)

		if context.InlineEdit {
			prefix := widgets.Resource.GetAdmin().GetRouter().Prefix

			return template.HTML(fmt.Sprintf(
				"<script data-prefix=\"%v\" src=\"%v/assets/javascripts/widget_check.js?theme=widget\"></script><div class=\"qor-widget qor-widget-%v\" data-widget-inline-edit-url=\"%v\" data-url=\"%v\">\n%v\n</div>",
				prefix,
				prefix,
				utils.ToParamString(widgetObj.Name),
				fmt.Sprintf("%v/%v/inline-edit", prefix, widgets.Resource.ToParam()),
				url,
				widgetObj.Render(widgetSetting, setting.GetTemplate(), url),
			))
		}

		return widgetObj.Render(widgetSetting, setting.GetTemplate(), url)
	}

	return template.HTML("")
}

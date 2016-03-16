package widget

import (
	"github.com/qor/qor/utils"
	"html/template"
)

func Render(name string, context *Context, availableWidgets ...string) template.HTML {
	if len(availableWidgets) == 0 {
		utils.ExitWithMsg("Widget Name can't be blank")
	}
	widgetName := availableWidgets[0]
	widgetObj, _ := GetWidget(widgetName)
	return widgetObj.Render()
}

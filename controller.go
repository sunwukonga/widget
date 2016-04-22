package widget

import (
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/responder"
)

type widgetController struct {
	Widgets *Widgets
}

func (wc widgetController) Index(context *admin.Context) {
	context = context.NewResourceContext(wc.Widgets.WidgetSettingResource)
	result, err := context.FindMany()
	context.AddError(err)

	if context.HasError() {
		http.NotFound(context.Writer, context.Request)
	} else {
		responder.With("html", func() {
			context.Execute("index", result)
		}).With("json", func() {
			context.JSON("index", result)
		}).Respond(context.Request)
	}
}

func (wc widgetController) Edit(context *admin.Context) {
	context.Resource = wc.Widgets.WidgetSettingResource
	widgetSetting := context.Resource.NewStruct()
	context.AddError(context.GetDB().First(widgetSetting, "name = ? AND scope = ?", context.ResourceID, "default").Error)
	context.Execute("edit", widgetSetting)
}

func (wc widgetController) Update(context *admin.Context) {
	context.Resource = wc.Widgets.WidgetSettingResource
	widgetSetting := context.Resource.NewStruct()
	context.AddError(context.GetDB().First(widgetSetting, "name = ? AND scope = ?", context.ResourceID, "default").Error)

	if context.AddError(context.Resource.Decode(context.Context, widgetSetting)); !context.HasError() {
		context.AddError(context.Resource.CallSave(widgetSetting, context.Context))
		context.Execute("edit", widgetSetting)
		return
	}

	http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
}

func (wc widgetController) InlineEdit(context *admin.Context) {
	context.Writer.Write([]byte(context.Render("inline_edit")))
}

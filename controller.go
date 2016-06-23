package widget

import (
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/responder"
	"github.com/qor/serializable_meta"
)

type widgetController struct {
	Widgets *Widgets
}

func (wc widgetController) Index(context *admin.Context) {
	context = context.NewResourceContext(wc.Widgets.WidgetSettingResource)
	result, _, err := wc.getWidget(context)
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
	widgetSetting, scopes, err := wc.getWidget(context)
	context.AddError(err)
	context.Execute("edit", map[string]interface{}{"Scopes": scopes, "Widget": widgetSetting})
}

func (wc widgetController) Update(context *admin.Context) {
	context.Resource = wc.Widgets.WidgetSettingResource
	widgetSetting, scopes, err := wc.getWidget(context)
	context.AddError(err)

	if context.AddError(context.Resource.Decode(context.Context, widgetSetting)); !context.HasError() {
		context.AddError(context.Resource.CallSave(widgetSetting, context.Context))
	}

	if context.HasError() {
		context.Writer.WriteHeader(admin.HTTPUnprocessableEntity)
		context.Execute("edit", map[string]interface{}{"Scopes": scopes, "Widget": widgetSetting})
	} else {
		http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
	}
}

func (wc widgetController) InlineEdit(context *admin.Context) {
	context.Writer.Write([]byte(context.Render("inline_edit")))
}

func (wc widgetController) getWidget(context *admin.Context) (interface{}, []string, error) {
	if context.ResourceID == "" {
		// index page
		context.SetDB(context.GetDB().Where("scope = ?", "default").Order("activated_at DESC").Group("name"))
		results, err := context.FindMany()
		return results, []string{}, err
	}

	// show page
	var (
		widgetSettings  []QorWidgetSetting
		selectedSetting *QorWidgetSetting
		scopes          []string
		result          = wc.Widgets.WidgetSettingResource.NewStruct()
		scope           = context.Request.URL.Query().Get("widget_scope")
		widgetType      = context.Request.URL.Query().Get("widget_type")
	)

	if scope == "" {
		scope = context.Request.Form.Get("QorResource.Scope")
	}

	if widgetType == "" {
		widgetType = context.Request.Form.Get("QorResource.Kind")
	}

	if scope == "" {
		scope = "default"
	}

	context.GetDB().Model(result).Where("name = ?", context.ResourceID).Order("activated_at").Find(&widgetSettings)

	for _, setting := range widgetSettings {
		if setting.Scope == scope {
			selectedSetting = &QorWidgetSetting{Name: context.ResourceID, Scope: scope}

			if setting.WidgetType == widgetType {
				selectedSetting = &QorWidgetSetting{Name: context.ResourceID, Scope: scope, WidgetType: widgetType}
				break
			}
		}
	}

	if selectedSetting == nil {
		selectedSetting = &QorWidgetSetting{Name: context.ResourceID, Scope: "default"}
	}

	err := context.GetDB().Order("activated_at DESC").First(result, selectedSetting).Error

	if widgetType != "" {
		if serializableMeta, ok := result.(serializable_meta.SerializableMetaInterface); ok && serializableMeta.GetSerializableArgumentKind() != widgetType {
			serializableMeta.SetSerializableArgumentKind(widgetType)
			serializableMeta.SetSerializableArgumentValue(nil)
		}
	}
	return result, scopes, err
}

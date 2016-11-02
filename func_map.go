package widget

import "github.com/qor/admin"

var funcMap = map[string]interface{}{
	"widget_available_scopes": func() []*Scope {
		if len(registeredScopes) > 0 {
			return append([]*Scope{{Name: "Default Visitor", Param: "default"}}, registeredScopes...)
		}
		return []*Scope{}
	},
	"get_widget_scopes": func(context *admin.Context) []string {
		_, scopes, _ := getWidget(context)
		return scopes
	},
}

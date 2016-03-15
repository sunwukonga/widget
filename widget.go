package widget

import "github.com/qor/admin"

type Widget struct {
	Name     string
	Template string
	Setting  *admin.Resource
	Context  func(context Context, setting interface{}) Context
}

var registeredWidgets []*Widget

func RegisterWidget(scope *Widget) {
	registeredWidgets = append(registeredWidgets, scope)
}

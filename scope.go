package widget

import "github.com/qor/qor/utils"

var registeredScopes []*Scope

// Scope widget scope
type Scope struct {
	Name    string
	Visible func(*Context) bool
}

func (scope *Scope) ToParam() string {
	return utils.ToParamString(scope.Name)
}

// RegisterScope register scope for widget
func (widgets *Widgets) RegisterScope(scope *Scope) {
	registeredScopes = append(registeredScopes, scope)
}

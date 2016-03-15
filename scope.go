package widget

type Scope struct {
	Name    string
	Visible func(*Context) bool
}

var registeredScopes []*Scope

func RegisterScope(scope *Scope) {
	registeredScopes = append(registeredScopes, scope)
}

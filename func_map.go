package widget

var funcMap = map[string]interface{}{
	"widget_available_scopes": func() (results []string) {
		if len(registeredScopes) > 0 {
			results = append(results, "Default")
		}

		for _, scope := range registeredScopes {
			results = append(results, scope.Name)
		}
		return
	},
}

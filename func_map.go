package widget

import "fmt"

func init() {
	var funcMap = map[string]interface{}{
		"widget_available_scopes": func() (results []string) {
			for _, scope := range registeredScopes {
				results = append(results, scope.Name)
			}
			return
		},
	}
	fmt.Println(funcMap)
}

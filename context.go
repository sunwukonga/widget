package widget

import "errors"

func NewContext(options map[string]interface{}, availableWidgets ...string) *Context {
	return &Context{
		Widgets: availableWidgets,
		Options: options,
	}
}

type Context struct {
	Widgets []string
	Options map[string]interface{}
}

func (context Context) Get(name string) (interface{}, error) {
	if value, ok := context.Options[name]; ok {
		return value, nil
	}

	return nil, errors.New("not found")
}

func (context *Context) Set(name string, value interface{}) {
	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}
	context.Options[name] = value
}

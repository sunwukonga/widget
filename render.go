package widget

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qor/qor/utils"
)

// Render find widget by name, render it based on current context
func (widgets *Widgets) Render(widgetName string, widgetGroupName string) template.HTML {
	return widgets.NewContext(nil).Render(widgetName, widgetGroupName)
}

// NewContext create new context for widgets
func (widgets *Widgets) NewContext(context *Context) *Context {
	if context == nil {
		context = &Context{}
	}

	if context.DB == nil {
		context.DB = widgets.Config.DB
	}

	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}

	if context.FuncMaps == nil {
		context.FuncMaps = template.FuncMap{}
	}

	for key, fc := range widgets.funcMaps {
		if _, ok := context.FuncMaps[key]; !ok {
			context.FuncMaps[key] = fc
		}
	}

	context.Widgets = widgets
	return context
}

// Funcs return view functions map
func (context *Context) Funcs(funcMaps template.FuncMap) *Context {
	if context.FuncMaps == nil {
		context.FuncMaps = template.FuncMap{}
	}

	for key, fc := range funcMaps {
		context.FuncMaps[key] = fc
	}

	return context
}

func (context *Context) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_widget"] = func(widgetName string, widgetGroupName ...string) template.HTML {
		var groupName string
		if len(widgetGroupName) == 0 {
			groupName = ""
		} else {
			groupName = widgetGroupName[0]
		}
		return context.Render(widgetName, groupName)
	}

	return funcMap
}

// Render register widget itself content
func (w *Widget) Render(context *Context, file string) template.HTML {
	var err error
	var result = bytes.NewBufferString("")
	if file == "" {
		file = w.Templates[0]
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Get error when render file %v: %v", file, r)
			utils.ExitWithMsg(err)
		}
	}()

	if file, err = w.findTemplate(file + ".tmpl"); err == nil {
		if tmpl, err := template.New(filepath.Base(file)).Funcs(context.FuncMaps).ParseFiles(file); err == nil {
			if err = tmpl.Execute(result, context.Options); err == nil {
				return template.HTML(result.String())
			}
		}
	}

	return template.HTML(err.Error())
}

// RegisterViewPath register views directory
func (widgets *Widgets) RegisterViewPath(p string) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		if registerViewPath(path.Join(gopath, "src", p)) == nil {
			return
		}
	}
}

// LoadPreviewAssets will return assets tag used for Preview
func (w *Widgets) LoadPreviewAssets() template.HTML {
	tags := ""
	for _, asset := range w.Config.PreviewAssets {
		extension := filepath.Ext(asset)
		if extension == ".css" {
			tags += fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%v\">\n", asset)
		} else if extension == ".js" {
			tags += fmt.Sprintf("<script src=\"%v\"></script>\n", asset)
		} else {
			tags += fmt.Sprintf("%v\n", asset)
		}
	}
	return template.HTML(tags)
}

func isExistingDir(pth string) bool {
	fi, err := os.Stat(pth)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func registerViewPath(path string) error {
	if isExistingDir(path) {
		var found bool

		for _, viewPath := range viewPaths {
			if path == viewPath {
				found = true
				break
			}
		}

		if !found {
			viewPaths = append(viewPaths, path)
		}
		return nil
	}
	return errors.New("path not found")
}

func (w *Widget) findTemplate(layouts ...string) (string, error) {
	for _, layout := range layouts {
		for _, p := range viewPaths {
			if _, err := os.Stat(filepath.Join(p, layout)); !os.IsNotExist(err) {
				return filepath.Join(p, layout), nil
			}
		}
	}
	return "", fmt.Errorf("template not found: %v", layouts)
}

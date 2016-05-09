# Widget

Create a widget for frontend

## Usage

### Getting Started

```go
// Initialize a new widgets container
Widgets := widget.New(&widget.Config{DB: db})

// Widget Settings Argument
type bannerArgument struct {
  Title           string
  Link            string
  BackgroundImage media_library.FileSystem
}

// Register a new widget
Widgets.RegisterWidget(&widget.Widget{
  Name:     "Banner",
  Templates: []string{"banner1", "banner2"},
  Setting:  Admin.NewResource(&bannerArgument{}),
  Context: func(context *widget.Context, setting interface{}) *widget.Context {
    context.Options["Setting"] = argument
    context.Options["CurrentTime"] = time.Now()
    return context
  },
})

// Add to qor admin
Admin.AddResource(Widgets)
```

### Templates

```go
// app/views/widgets/banner1.tmpl
<div class="banner" style="background:url('{{.Setting.BackgroundImage}}') no-repeat center center">
  <div class="container">
    <div class="row">
      <div class="column column-12">
        <a href="{{.Setting.Link}}" class="button button__primary">{{.Setting.Title}}</a>
        {{.CurrentTime}}
      </div>
    </div>
  </div>
</div>

// app/views/widgets/banner2.tmpl
<div class="banner">
  <div class="column column-9">
    <img src="{{.Setting.BackgroundImage}}" class="banner__logo" />
  </div>

  <div class="column column-3" style="margin-top: 2em;">
    <div><a href="{{.Setting.Link}}" class="button button__primary">{{.Setting.Title}}</a></div>
  </div>
</div>
```

### Register Scopes

```go
Widgets.RegisterScope(&widget.Scope{
	Name: "From Google",
	Visible: func(context *widget.Context) bool {
		if request, ok := context.Get("Request"); ok {
			return strings.Contains(request.(*http.Request).Referer(), "google.com")
		}
		return false
	},
})

Widgets.RegisterScope(&widget.Scope{
	Name: "VIP User",
	Visible: func(context *widget.Context) bool {
		if user, ok := context.Get("CurrentUser"); ok {
			return user.(*User).Role == "vip"
		}
		return false
	},
})
```

### Render Widget

```go
func Index(request *http.Request, writer http.ResponseWriter) {
  widgetContext := widget.NewContext(map[string]interface{}{"Request": request, "CurrentUser": currentUser})

  // Render Widget
  bannerContent := Widgets.Render("Banner", "HomeBanner",  widgetContext)

  // Render Widget With Inline Edit
  bannerContent := Widgets.Render("Banner", "HomeBanner",  widgetContext, true)
}
```

## Live Demo

[Qor Demo](http://demo.getqor.com)

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).

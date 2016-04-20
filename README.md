# Widget

Web Widgets

### Configuration

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
  Templates: []string{"banner", "slideout"},
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
// app/views/widgets/banner.tmpl
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
```

### Render

```go
Widgets.Render("Banner", "top_banner", nil)
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).

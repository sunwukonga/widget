# Layout

Page Builder - WIP

### Configuration

```go
RegisterViewPath "templates"

RegisterScope(layout.Scope{
  Name: "From Google", // User purchased in last month
  Visible: func() bool {
    // return
  },
})

RegisterWidget(layout.Widget{
  Name:  "Home Header",
  Requires: []{"Mini Cart", "navigation"},
  Template: "home_header",
  Setting: *admin.Resource,
  Context: func(Context) map[string]interface{} {
  },
})

RegisterWidget(layout.Widget{
  Name:  "Mini Cart",
  Template: "mini_cart",
  Setting: *admin.Resource,
  Context: func(Context) map[string]interface{} {
  },
})

context := layout.NewContext(Context{
  AvailableWidget: []string{"Home Header"},
  Options: map[string]interface{}{
    "CurrentUser": user,
    "CurrentProduct": product,
  }},
)

func (Context) Get(string) interface{} {
}
```

### Page

```html
{{render_qor_layout "Qor Home Header" context}}
```

### Template

```go
<div>
  <div class="col-lg-4">
    {{embed_qor_layout "logo" "logo"}}
  </div>

  <div class="col-lg-8">
     {{.SearchURL}}
  </div>

  <div class="col-lg-12">
    {{render_qor_layout "Mini Cart"}}
  </div>
</div>
```

### Database

```csv
// qor_layouts
name, scope, widget_name, settings
"home_page_header", "form google", "widget_name", {}
"logo", "from_baidu", "widget_name", {}
```

### Func Map

```
render_qor_layout -> get scope name, search with scope & layout name, get widget name, decode argument, render
```

# Layout

Page Builder - WIP

### Configuration

```go
RegisterViewPath "templates"

RegisterWidget(layout.Widget{
  Name:  "Banner",
  Requires: []{"Mini Cart", "navigation"},
  Template: "home_header",
  Setting: *admin.Resource,
  Context: func(context Context, setting interface{}) Context {
   "widget"
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
  AvailableWidgets: []string{"Home Header"},
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
{{render_widget "Qor Home Header" context}}

{{render_widget "Product Show" context}}
  - {{render_widget "Cart" context}}
```

### Template

```go
<div>
  <div class="col-lg-4">
    {{embed_widget "logo" "logo"}}
  </div>

  <div class="col-lg-8">
     {{.SearchURL}}
  </div>

  <div class="col-lg-12">
    {{render_widget "Mini Cart"}}
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

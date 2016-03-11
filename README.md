# Layout

Page Builder - WIP

### Configuration

```go
RegisterViewPath "templates"

RegisterWidget(layout.Widget{
  Name:  "Home Header",
  Template: "home_header",
  Setting: *admin.Resource,
  Context: func() map[string]interface{} {
  },
})
```

### Page

```html
{{render_qor_layout "home_page_header"}}
```

### Template

```go
<div>
  <div class="col-lg-4">
    {{embed_qor_layout "logo"}}
  </div>

  <div class="col-lg-8">
     {{.SearchURL}}
  </div>
</div>
```

### Database

```csv
// qor_layouts
name, kind, settings
"home_page_header", "widget_name", {}
"logo", "widget_name", {}
```

### Func Map

```
render_qor_layout -> search with name, get widget name, decode argument, render
```

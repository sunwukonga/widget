package widget_test

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/test/utils"
	"github.com/qor/widget"
)

var db *gorm.DB

func init() {
	db = utils.TestDB()
}

// Runner
func TestRender(t *testing.T) {
	widget.RegisterViewPath("github.com/qor/widget/test")
	Admin := admin.New(&qor.Config{DB: db})

	type bannerArgument struct {
		Title    string
		SubTitle string
	}

	widget.RegisterWidget(&widget.Widget{
		Name:     "Banner",
		Template: "banner",
		Setting:  Admin.NewResource(&bannerArgument{}),
		Context: func(context widget.Context, setting interface{}) widget.Context {
			return context
		},
	})

	widgetContext := widget.NewContext(map[string]interface{}{
		"Title":    "Hello world!",
		"SubTitle": "Banner Widget",
	})
	html := widget.Render("HomeBanner", *widgetContext, "Banner")
	if html != "<h1>Hello world!</h1>\n<h2>Banner Widget</h2>\n" {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result: %s\n", 1, html)))
	}
}

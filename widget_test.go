package widget_test

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/test/utils"
	"github.com/qor/widget"
)

var db *gorm.DB

func init() {
	db = utils.TestDB()
}

// Runner
func TestRender(t *testing.T) {
	widget.RegisterWidget(&widget.Widget{
		Name:     "Banner",
		Template: "banner.tmpl",
		Context:  func(context widget.Context, setting interface{}) widget.Context { return widget.Context{} },
	})

	widgetContext := widget.NewContext(map[string]interface{}{})
	html := widget.Render("HomeBanner", widgetContext, "Banner")
	if html != "Hello" {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result: %s\n", 1, html)))
	}
}

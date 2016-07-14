package widget_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/test/utils"
	"github.com/qor/widget"
)

var db *gorm.DB
var Widgets *widget.Widgets
var Admin *admin.Admin
var Server *httptest.Server

type bannerArgument struct {
	Title    string
	SubTitle string
}

func init() {
	db = utils.TestDB()
	mux := http.NewServeMux()
	Server = httptest.NewServer(mux)

	Widgets = widget.New(&widget.Config{
		DB: db,
	})
	Widgets.RegisterViewPath("github.com/qor/widget/test")

	Admin = admin.New(&qor.Config{DB: db})
	Admin.AddResource(Widgets)
	Admin.MountTo("/admin", mux)

	Widgets.RegisterWidget(&widget.Widget{
		Name:      "Banner",
		Templates: []string{"banner"},
		Setting:   Admin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			if setting != nil {
				argument := setting.(*bannerArgument)
				context.Options["Title"] = argument.Title
				context.Options["SubTitle"] = argument.SubTitle
			}
			return context
		},
	})

	Widgets.RegisterScope(&widget.Scope{
		Name: "From Google",
		Visible: func(context *widget.Context) bool {
			if request, ok := context.Get("Request"); ok {
				_, ok := request.(*http.Request).URL.Query()["from_google"]
				return ok
			}
			return false
		},
	})
}

func reset() {
	db.DropTable(&widget.QorWidgetSetting{})
	db.AutoMigrate(&widget.QorWidgetSetting{})
}

// Test DB's record after call Render
func TestRenderRecord(t *testing.T) {
	reset()
	var count int
	db.Model(&widget.QorWidgetSetting{}).Where(widget.QorWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "default", GroupName: "Banner"}).Count(&count)
	if count != 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should don't exist widget setting")))
	}

	widgetContext := Widgets.NewContext(&widget.Context{})
	widgetContext.Render("HomeBanner", "Banner")
	db.Model(&widget.QorWidgetSetting{}).Where(widget.QorWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "default", GroupName: "Banner"}).Count(&count)
	if count == 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should have default widget setting")))
	}

	http.PostForm(Server.URL+"/admin/widgets/HomeBanner",
		url.Values{"_method": {"PUT"},
			"QorResource.Scope":       {"from_google"},
			"QorResource.ActivatedAt": {"2016-07-14 10:10:42.433372925 +0800 CST"},
			"QorResource.Widgets":     {"Banner"},
			"QorResource.Template":    {"banner"},
			"QorResource.Kind":        {"Banner"},
		})
	db.Model(&widget.QorWidgetSetting{}).Where(widget.QorWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "from_google"}).Count(&count)
	if count == 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should have from_google widget setting")))
	}
}

// Runner
func TestRenderContext(t *testing.T) {
	reset()
	setting := &widget.QorWidgetSetting{}
	db.Where(widget.QorWidgetSetting{Name: "HomeBanner", WidgetType: "Banner"}).FirstOrInit(setting)
	db.Create(setting)

	html := Widgets.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, \n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 1, html)))
	}

	widgetContext := Widgets.NewContext(&widget.Context{
		Options: map[string]interface{}{"CurrentUser": "Qortex"},
	})
	html = widgetContext.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, Qortex\n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 2, html)))
	}

	db.Where(widget.QorWidgetSetting{Name: "HomeBanner", WidgetType: "Banner"}).FirstOrInit(setting)
	setting.SetSerializableArgumentValue(&bannerArgument{Title: "Title", SubTitle: "SubTitle"})
	db.Save(setting)

	html = widgetContext.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, Qortex\n<h1>Title</h1>\n<h2>SubTitle</h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 3, html)))
	}
}

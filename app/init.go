package app

import (
	"log"
	"myapp/app/utils"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

var (
	Engine *xorm.Engine
)

const (
	CaptchaW = 120
	CaptchaH = 40
)

func init() {

	// init revel filters
	initRevelFilters()
	// init revel template functions
	initRevelTemplateFuncs()

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(initHandlers)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

func initRevelFilters() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}
}

func initRevelTemplateFuncs() {
	revel.TemplateFuncs["webTitle"] = func(prefix string) (webTitle string) {
		const KEY = "cache.web.title"
		if err := cache.Get(KEY, &webTitle); err != nil {
			webTitle = utils.ForceGetConfig("web.title")
			go cache.Set(KEY, webTitle, 24*30*time.Hour)
		}
		return
	}

	revel.TemplateFuncs["logined"] = func(session revel.Session) bool {
		v, e := session["user"]
		return e == true && v != ""
	}

	revel.TemplateFuncs["isAdmin"] = func(session revel.Session) bool {
		return false
	}
}

func initHandlers() {
	var (
		serveMux     = http.NewServeMux()
		revelHandler = revel.Server.Handler
	)
	serveMux.Handle("/", revelHandler)
	// serveMux.Handle("/captcha/", http.HandlerFunc(CaptchaHanlder))
	serveMux.Handle("/captcha/", captcha.Server(CaptchaW, CaptchaH))
	revel.Server.Handler = serveMux

}

var CaptchaHanlder = func(w http.ResponseWriter, r *http.Request) {
	log.Println("CaptchaHandler handling...")
	captcha.Server(CaptchaW, CaptchaH).ServeHTTP(w, r)
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

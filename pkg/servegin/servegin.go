package servegin

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/pathtraversal"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var base = gin.H{"Framework": "Gin", "Logo": "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png"}

//return a copy of 'base', with Name updated. uses a copy to avoid race conditions.
func templateData(name string) gin.H {
	cpy := make(gin.H, len(base)+1)
	for k, v := range base {
		cpy[k] = v
	}
	cpy["Name"] = name
	return cpy
}

func loadTemplates() multitemplate.Renderer {
	templatesDir, err := common.FindViewsDir()
	if err != nil {
		panic(err.Error())
	}
	pages, err := filepath.Glob(filepath.Join(templatesDir, "pages", "*.gohtml"))
	if err != nil {
		panic(err.Error())
	}
	partials, err := filepath.Glob(filepath.Join(templatesDir, "partials", "*.gohtml"))
	if err != nil {
		panic(err.Error())
	}
	layout := filepath.Join(templatesDir, "layout.gohtml")

	r := multitemplate.NewRenderer()
	for _, p := range pages {
		files := append([]string{layout, p}, partials...)
		r.AddFromFilesFuncs(filepath.Base(p), common.FuncMap(), files...)
	}

	return r
}

//add a handler to gin
func add(router *gin.Engine, rt common.Route) {
	base := router.Group(rt.Base)
	base.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, rt.TmplFile, templateData(rt.Base))
	})
	for _, s := range rt.Sinks {
		sinkFn := func(s common.Sink) func(c *gin.Context) {
			return func(c *gin.Context) {
				c.Header("Cache-Control", "no-store") //makes development a whole lot easier
				mode := common.Safety(c.Param("mode"))
				source := c.Param("source")
				payload := extractInput(c, source)
				for _, e := range c.Errors {
					log.Printf("%s: error %s", c.Request.URL.Path, e)
				}

				var data string
				if s.Handler != nil {
					data = string(s.Handler(mode, payload, c))
				} else {
					data = string(common.GenericHandler(s, mode, payload, c))
				}
				if len(data) > 0 {
					// don't unconditionally write this, as it can result in
					// - a warning (when status changes), or
					// - a panic (when content-length is already set and headers are written)
					c.String(http.StatusOK, data)
				}
			}
		}(s)
		sinkPg := base.Group("/" + s.URL)
		//route data isn't a perfect match for the method(s) we actually use, so just accept anything
		sinkPg.Any("/:source/:mode", sinkFn)
	}
}

var ginPathTraversal = common.Sink{
	Name:     "gin.File",
	Sanitize: url.QueryEscape,
	VulnerableFnWrapper: func(opaque interface{}, payload string) (data string, err error) {
		c, ok := opaque.(*gin.Context)
		if !ok {
			log.Fatalf("'opaque': want *gin.Context, got %T", opaque)
		}
		c.File(payload)
		return "", common.NoDecoration
	},
}

// Setup loads templates, sets up routes, etc.
func Setup(addr string) (router *gin.Engine, dbFile string) {
	base["Addr"] = addr

	//register all routes at this point, before AllRoutes is used.
	cmdi.RegisterRoutes("gin")
	sqli.RegisterRoutes("gin")
	pathtraversal.RegisterRoutes([]common.Sink{ginPathTraversal})

	rmap := common.PopulateRouteMap(common.AllRoutes)

	//until all routes are migrated to the new model, we need to do a few fixups
	rmap = preMigrationFixups(rmap)

	base["Rulebar"] = rmap
	router = gin.Default()

	log.Println("Loading templates...")
	router.HTMLRender = loadTemplates()
	log.Println("Templates loaded.")

	router.StaticFS("/assets/", http.Dir("./public"))
	router.GET("/", func(c *gin.Context) {
		c.Header("Application-Framework", "Gin")
		c.HTML(http.StatusOK, "index.gohtml", templateData(""))
	})

	for _, h := range common.AllRoutes {
		add(router, h)
	}
	addReflectedXSS(router)
	addSSRF(router)
	addUnvalidatedRedirect(router)

	// setting up a database to execute the built query
	dbSrc, err := os.CreateTemp(".", "tempDatabase*.db")
	if err != nil {
		panic(err)
	}
	return router, dbSrc.Name()
}

//temporary fixes until remainder of code migrates to new model
func preMigrationFixups(rmap common.RouteMap) common.RouteMap {
	// unvalidated redirect; for now, just handle the gin method
	ur, ok := rmap["unvalidatedRedirect"]
	if !ok {
		for k := range rmap {
			log.Println(k)
		}
		log.Fatal("unvalidated redirect is missing")
	}
	ur.Sinks[0].Name = "gin.Redirect"
	ur.Sinks[0].URL = strings.Replace(ur.Sinks[0].URL, "http.", "gin.", 1)
	rmap["unvalidatedRedirect"] = ur

	return rmap
}

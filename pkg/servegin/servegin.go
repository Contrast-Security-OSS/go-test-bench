package servegin

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
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

	fmap := template.FuncMap{"tolower": strings.ToLower}

	r := multitemplate.NewRenderer()
	for _, p := range pages {
		files := append([]string{layout, p}, partials...)
		r.AddFromFilesFuncs(filepath.Base(p), fmap, files...)
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
		sinkFn := func(c *gin.Context) {
			mode := c.Param("mode")
			source := c.Param("source")
			payload := extractInput(c, source)

			tmpl, b := s.Handler(mode, payload)
			if b {
				log.Fatal("error: bool arg is not handled")
			}
			c.String(http.StatusOK, string(tmpl))
		}
		sinkPg := base.Group("/" + s.URL)
		//route data isn't a perfect match for the method(s) we actually use, so just accept anything
		sinkPg.Any("/:source/:mode", sinkFn)
	}
}

// Setup loads templates, sets up routes, etc.
func Setup(addr string) (router *gin.Engine, dbFile string) {
	base["Addr"] = addr

	//register all routes at this point, before AllRoutes is used.
	cmdi.RegisterRoutes("gin")

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
	addPathTraversal(router)
	addReflectedXSS(router)
	addSSRF(router)
	addUnvalidatedRedirect(router)

	// setting up a database to execute the built query
	dbSrc, err := os.CreateTemp(".", "tempDatabase*.db")
	if err != nil {
		panic(err)
	}
	addSQLi(router, dbSrc)

	return router, dbSrc.Name()
}

//temporary fixes until remainder of code migrates to new model
func preMigrationFixups(rmap common.RouteMap) common.RouteMap {
	//for path traversal, gin supports an additional method
	//this will go into path traversal's RegisterRoutes() func when it's migrated
	pt, ok := rmap["pathTraversal"]
	if !ok {
		for k := range rmap {
			log.Println(k)
		}
		log.Fatal("path traversal is missing")
	}
	pt.Sinks = append(
		[]common.Sink{{
			Name:   "gin.File",
			URL:    "/pathTraversal",
			Method: "GET",
		}},
		pt.Sinks...,
	)
	rmap["pathTraversal"] = pt

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

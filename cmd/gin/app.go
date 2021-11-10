package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

var base = gin.H{"Framework": "Gin", "Logo": "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png"}

//return a copy of 'base', with Name updated. uses a copy to avoid race conditions.
func templateData(name string) gin.H {
	cpy := make(gin.H)
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
		sinkPg.GET("/:source/:mode", sinkFn)
		sinkPg.POST("/:source/:mode", sinkFn)
	}
}

func main() {
	// Setup command line flags
	port := flag.Int("port", DefaultPort, "listen on this `port` on localhost")
	flag.Parse()
	addr := fmt.Sprintf("localhost:%d", *port)
	base["Addr"] = addr

	//register all routes at this point, before AllRoutes is used.
	cmdi.RegisterRoutes("gin")

	base["Rulebar"] = common.PopulateRouteMap(common.AllRoutes)

	router := gin.Default()

	log.Println("Loading templates...")
	router.HTMLRender = loadTemplates()
	log.Println("Templates loaded.")

	router.StaticFS("/assets/", http.Dir("./public"))
	router.GET("/", func(c *gin.Context) {
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

	// graceful shutdown to clean up database file
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Shutting down")
		err := os.Remove(dbSrc.Name())
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	log.Printf("Server startup at: %s\n", addr)
	log.Fatal(router.Run(addr))
}

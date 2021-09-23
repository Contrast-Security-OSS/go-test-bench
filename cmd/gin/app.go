package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

var base = gin.H{"Rulebar": rules, "Framework": "Gin", "Logo": "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png"}

func templateData(name string) gin.H {
	base["Name"] = name
	return base
}

func loadTemplates() multitemplate.Renderer {
	templatesDir := filepath.Clean("./views")
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
		r.AddFromFiles(filepath.Base(p), files...)
	}

	return r
}

func main() {
	// Setup command line flags
	portPtr := flag.Int("port", DefaultPort, "listen on this port")
	flag.Parse()
	portAddr := fmt.Sprintf(":%d", *portPtr)
	base["Port"] = portAddr

	router := gin.Default()

	log.Println("Loading templates...")
	router.HTMLRender = loadTemplates()
	log.Println("Templates loaded.")

	router.StaticFS("/assets/", http.Dir("./public/gin"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", templateData(""))
	})

	addCMDi(router)
	addPathTraversal(router)
	addReflectedXSS(router)
	addSSRF(router)
	addUnvalidatedRedirect(router)

	// setting up a database to execute the built query
	dbSrc, err := os.CreateTemp(".", "tempDatabase*.db")
	if err != nil {
		panic(err)
	}
	defer os.Remove(dbSrc.Name())
	addSQLi(router, dbSrc)

	log.Printf("Server startup at: localhost%s\n", portAddr)
	log.Fatal(router.Run(portAddr))
}

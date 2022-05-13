package servegin

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func addReflectedXSS(r *gin.Engine) {
	xss := r.Group("/xss")
	xss.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "xss.gohtml", templateData("xss"))
	})

	xss.GET("/:source/reflectedXss/:type", xssHandlerFunc)
	xss.POST("/:source/reflectedXss/:type", xssHandlerFunc)
	xss.GET("/:source/reflectedXss/:type/*param", xssHandlerFunc)
}

func xssHandlerFunc(c *gin.Context) {
	source := c.Param("source")
	payload := extractInput(c, source)

	switch c.Param("type") {
	case "noop":
		payload = "noop"
	case "safe":
		payload = url.QueryEscape(payload)
	case "unsafe":
	}

	c.Writer.WriteHeader(http.StatusOK)
	_, err := c.Writer.WriteString(payload)
	if err != nil {
		log.Printf("writing xss payload %s for %s: error %s", payload, c.Request.URL.Path, err)
	}
}

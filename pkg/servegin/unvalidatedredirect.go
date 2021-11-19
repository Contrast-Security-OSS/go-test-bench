package servegin

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func addUnvalidatedRedirect(r *gin.Engine) {
	ur := r.Group("/unvalidatedRedirect")
	ur.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "unvalidatedRedirect.gohtml", templateData("unvalidatedRedirect"))
	})

	ur.GET("/gin.Redirect/:type", unvalidatedRedirectHandlerFunc)
}

func unvalidatedRedirectHandlerFunc(c *gin.Context) {
	payload := extractInput(c, "query")

	switch c.Param("type") {
	case "noop":
		c.String(http.StatusOK, "noop")
		return
	case "safe":
		payload = url.PathEscape(payload)
	case "unsafe":
	}

	c.Redirect(http.StatusTemporaryRedirect, payload)
}

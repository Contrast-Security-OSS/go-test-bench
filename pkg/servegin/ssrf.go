package servegin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addSSRF(r *gin.Engine) {
	ssrf := r.Group("/ssrf")
	ssrf.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ssrf.gohtml", templateData("ssrf"))
	})

	ssrf.GET("/:lib/:vector", ssrfHandlerFunc)
}

func ssrfHandlerFunc(c *gin.Context) {
	payload := extractInput(c, "query")

	url := "http://example.com"
	switch c.Param("vector") {
	case "query":
		if payload != "" {
			url += "?input=" + payload
		}
	case "path":
		if payload != "" {
			url = "http://" + payload
		}
	}

	res, err := http.Get(url)
	if err != nil {
		c.Error(err)
		return
	}

	c.DataFromReader(
		http.StatusOK,
		res.ContentLength,
		"text/html",
		res.Body,
		nil,
	)
}

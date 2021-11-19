package servegin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func addPathTraversal(r *gin.Engine) {
	pt := r.Group("/pathTraversal")
	pt.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pathTraversal.gohtml", templateData("pathTraversal"))
	})

	pt.GET("/:source/:rw/:type", ptHandlerFunc)
	pt.POST("/:source/:rw/:type", ptHandlerFunc)
}

func ptHandlerFunc(c *gin.Context) {
	source := c.Param("source")
	payload := extractInput(c, source)

	switch c.Param("type") {
	case "noop":
		c.String(http.StatusOK, "noop")
		return
	case "safe":
		payload = url.QueryEscape(payload)
	case "unsafe":
	}

	var out string
	switch c.Param("rw") {
	case "gin.File":
		if _, err := os.Stat(payload); err != nil {
			// disambiguate - the 404 is for the file, not for an incorrect endpoint
			fmt.Fprintf(c.Writer, "endpoint hit; gin.File returns 404 for not found\n")
		}
		c.File(payload)
		return
	case "ioutil.ReadFile":
		data, err := ioutil.ReadFile(payload)
		if err != nil {
			c.Error(err)
			return
		}
		out = fmt.Sprintf("Successfully read file %s. Length: %d", payload, len(data))
	case "ioutil.WriteFile":
		err := ioutil.WriteFile(payload, []byte("haxx'd"), 0666)
		if err != nil {
			c.Error(err)
			return
		}
		out = fmt.Sprintf("Successfully wrote file %s.", payload)
	}

	c.String(http.StatusOK, out)
}

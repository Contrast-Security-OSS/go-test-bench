package servegin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func logErr(c *gin.Context, err error, str ...string) {
	if err == nil {
		return
	}
	log.Printf("getting input from %s: %s", c.Request.URL.Path, err)
	if len(str) > 0 {
		log.Println("vals:", str)
	}
	_ = c.Error(err)
}

func extractInput(c *gin.Context, source string) string {
	switch source {
	case "query":
		return c.Query("input")
	case "buffered-query":
		return buffer(c.Query("input"))
	case "params":
		// path parameter includes leading slash, so we chop it off.
		return c.Param("param")[1:]
	case "body":
		return c.PostForm("input")
	case "buffered-body":
		input := c.PostForm("input")
		return buffer(input)
	case "cookies":
		input, err := c.Cookie("input")
		logErr(c, err)
		return input
	case "headers":
		return c.GetHeader("input")
	case "headers-json":
		// currently only used for SQLi
		input := c.GetHeader("credentials")
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.Unmarshal([]byte(input), &creds)
		logErr(c, err, input)
		return creds.Username

	default:
		logErr(c, fmt.Errorf("invalid source: %s", source))
		return ""
	}
}

func buffer(s string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(s)
	return buf.String()
}

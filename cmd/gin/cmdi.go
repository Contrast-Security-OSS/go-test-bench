package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func addCMDi(r *gin.Engine) {
	cmdi := r.Group("/cmdInjection")
	cmdi.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "commandInjection.gohtml", templateData("cmdInjection"))
	})

	osExec := cmdi.Group("/osExec")
	osExec.GET("/:source/:type", cmdiHandlerFunc)
	osExec.POST("/:source/:type", cmdiHandlerFunc)
}

func cmdiHandlerFunc(c *gin.Context) {
	source := c.Param("source")
	payload := extractInput(c, source)

	var cmd *exec.Cmd
	switch c.Param("type") {
	case "noop":
		c.String(http.StatusOK, "noop")
		return
	case "safe":
		cmd = exec.Command("echo", payload)
	case "unsafe":
		cmd = exec.Command(payload)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not run command %q. err: %v", cmd, err)
	}

	c.String(http.StatusOK, out.String())
}

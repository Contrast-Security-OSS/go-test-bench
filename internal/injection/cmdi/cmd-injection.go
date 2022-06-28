package cmdi

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := []*common.Sink{
		{
			Name:    "exec.Command",
			Handler: execHandler,
		},
		{
			Name:    "exec.CommandContext",
			Handler: execHandlerCtx,
		},
	}
	sinks = append(sinks, frameworkSinks...)
	common.Register(common.Route{
		Name:     "Command Injection",
		Link:     "https://www.owasp.org/index.php/Command_Injection",
		Base:     "cmdInjection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "cookies"},
		Sinks:    sinks,
		Payload:  "hello there! && echo hack hack hack",
	})
}

const mime = "text/plain"

// perform the vulnerability, using exec.Command
func execHandler(mode common.Safety, in string, _ interface{}) (string, string, int) {

	var cmd *exec.Cmd
	switch mode {
	case common.Safe:
		cmd = exec.Command("echo", in)
	case common.Unsafe:
		args := shellArgs(in)
		if len(args) == 0 {
			return "one or more args required", mime, http.StatusBadRequest
		}
		cmd = exec.Command(args[0], args[1:]...)
	case common.NOOP:
		return "NOOP", mime, http.StatusOK
	default:
		log.Fatalf("Error running execHandler. Unknown option  %q passed", mode)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("Could not run command in %s: err = %s", mode, err)
		if _, err = out.WriteString(msg); err != nil {
			log.Print("failed to add error to returned data:", err)
		}
		log.Print(msg)
	}
	return out.String(), mime, http.StatusOK
}

// perform the vulnerability, using exec.CommandContext
func execHandlerCtx(mode common.Safety, in string, _ interface{}) (string, string, int) {
	var cmd *exec.Cmd
	ctx := context.Background()
	switch mode {
	case common.Safe:
		cmd = exec.CommandContext(ctx, "echo", in)
	case common.Unsafe:
		args := shellArgs(in)
		if len(args) == 0 {
			return "one or more args required", mime, http.StatusBadRequest
		}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	case common.NOOP:
		return "NOOP", mime, http.StatusOK
	default:
		log.Fatalf("Error running execHandlerCtx. Unknown option  %q passed", mode)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("Could not run command in %s: err = %s", mode, err)
		out.WriteString(msg)
		log.Print(msg)
	}
	return out.String(), mime, http.StatusOK
}

// assembles a command that will run unsanitized user input in a system shell
func shellArgs(in string) []string {
	return strings.Fields(in)
}

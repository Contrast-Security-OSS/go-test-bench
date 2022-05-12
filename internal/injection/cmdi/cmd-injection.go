package cmdi

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes( /* framework - unused */ string) {
	common.Register(common.Route{
		Name:     "Command Injection",
		Link:     "https://www.owasp.org/index.php/Command_Injection",
		Base:     "cmdInjection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "cookies"},
		Sinks: []common.Sink{
			{
				Name:    "exec.Command",
				Method:  "GET",
				Handler: execHandler,
			},
			{
				Name:    "exec.CommandContext",
				Method:  "GET",
				Handler: execHandlerCtx,
			},
		},
		Payload: "hello there! && echo hack hack hack",
	})
}

// perform the vulnerability, using exec.Command
func execHandler(mode common.Safety, in string, _ interface{}) template.HTML {
	var cmd *exec.Cmd
	switch mode {
	case common.Safe:
		cmd = exec.Command("echo", in)
	case common.Unsafe:
		args := shellArgs(in)
		if len(args) == 0 {
			break
		}
		cmd = exec.Command(args[0], args[1:]...)
	case common.NOOP:
		return template.HTML("NOOP")
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
	return template.HTML(out.String())
}

// perform the vulnerability, using exec.CommandContext
func execHandlerCtx(mode common.Safety, in string, _ interface{}) template.HTML {
	var cmd *exec.Cmd
	ctx := context.Background()
	switch mode {
	case common.Safe:
		cmd = exec.CommandContext(ctx, "echo", in)
	case common.Unsafe:
		args := shellArgs(in)
		if len(args) == 0 {
			break
		}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	case common.NOOP:
		return template.HTML("NOOP")
	default:
		log.Fatalf("Error running execHandlerCtx. Unknown option  %q passed", mode)
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
	return template.HTML(out.String())
}

// assembles a command that will run unsanitized user input in a system shell
func shellArgs(in string) []string {
	return strings.Fields(in)
}

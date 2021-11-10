package cmdi

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os/exec"
	"runtime"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

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
	})
}

// perform the vulnerability, using exec.Command
func execHandler(mode, in string) (template.HTML, bool) {
	var cmd *exec.Cmd
	switch mode {
	case "safe":
		cmd = exec.Command("echo", in)
	case "unsafe":
		args := shellArgs(in)
		cmd = exec.Command(args[0], args[1:]...)
	case "noop":
		return template.HTML("NOOP"), false
	default:
		log.Fatalf("Error running execHandler. Unknown option  %q passed", mode)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not run command in %s: err = %s", mode, err)
	}
	return template.HTML(out.String()), false
}

// perform the vulnerability, using exec.CommandContext
func execHandlerCtx(mode, in string) (template.HTML, bool) {
	var cmd *exec.Cmd
	ctx := context.Background()
	switch mode {
	case "safe":
		cmd = exec.CommandContext(ctx, "echo", in)
	case "unsafe":
		args := shellArgs(in)
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	case "noop":
		return template.HTML("NOOP"), false
	default:
		log.Fatalf("Error running execHandlerCtx. Unknown option  %q passed", mode)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not run command in %s: err = %s", mode, err)
	}
	return template.HTML(out.String()), false
}

// assembles a command that will run unsanitized user input in a system shell
func shellArgs(in string) []string {
	var args []string
	if runtime.GOOS == "windows" {
		ps, err := exec.LookPath("powershell.exe")
		if err == nil {
			args = append(args, ps, in)
		}
	}
	if len(args) == 0 {
		args = append(args, "sh", "-c", in)
	}
	return args
}

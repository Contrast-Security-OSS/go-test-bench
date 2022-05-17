package serveswagger

import (
	"log"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// CmdInjectionFront serves the front page of the command injection page
func CmdInjectionFront(params cmd_injection.CmdInjectionFrontParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		cmdInjectionRoot := "/cmdInjection"

		if err := p.Produce(w, cmdInjectionRoot); err != nil {
			log.Print(err.Error())
		}
	})

}

// GetQueryCommand passes the input from user query and then calls into cmdi.ExecHandler to perform a command injection
func GetQueryCommand(params cmd_injection.GetQueryCommandParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		var payload string

		txt, isTemplate := cmdi.ExecHandler(params.Safety, params.Input)
		if !isTemplate {
			payload = string(txt)
		}

		if err := p.Produce(w, payload); err != nil {
			log.Print(err.Error())
		}
	})
}

// GetQueryCommandContext passes the input from user query and then calls into cmdi.ExecHandleCtx to perform a command injection
func GetQueryCommandContext(params cmd_injection.GetQueryCommandContextParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		var payload string

		txt, isTemplate := cmdi.ExecHandlerCtx(params.Safety, params.Input)
		if !isTemplate {
			payload = string(txt)
		}

		if err := p.Produce(w, payload); err != nil {
			log.Print(err.Error())
		}
	})
}

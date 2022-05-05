package serveswagger

import (
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func CommandInjectionHandler(params cmd_injection.CmdInjectionFrontParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {

		if err := p.Produce(w, "/cmdInjection"); err != nil {
		}
	})

}
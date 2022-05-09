package main

import (
	"log"
	"net/http"
	"github.com/go-openapi/runtime"
	"os"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-openapi/loads"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	if err := serveswagger.Setup(); err != nil {
		log.Fatal(err)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewSwaggerBenchAPI(swaggerSpec)

	api.SwaggerServerRootHandler = swagger_server.RootHandlerFunc(serveswagger.SwaggerRootHandler)

	api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(serveswagger.CommandInjectionHandler)

	api.CmdInjectionGetQueryCommandHandler = cmd_injection.GetQueryCommandHandlerFunc(func(params cmd_injection.GetQueryCommandParams) middleware.Responder {
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			var data string

			txt, isTemplate := cmdi.ExecHandlerCtx(params.Safety, params.Input)
			if !isTemplate {
				data = string(txt)
			} else {
				data = "not intended to render"
			}

			if err := p.Produce(w, data); err != nil {
			}
		})
	})

	api.CmdInjectionGetQueryCommandContextHandler = cmd_injection.GetQueryCommandContextHandlerFunc(func(params cmd_injection.GetQueryCommandContextParams) middleware.Responder {
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			var data string

			txt, isTemplate := cmdi.ExecHandlerCtx(params.Safety, params.Input)
			if !isTemplate {
				data = string(txt)
			} else {
				data = "not intended to render"
			}

			if err := p.Produce(w, data); err != nil {
			}
		})
	})

	api.CmdInjectionPostCookiesCommandHandler = cmd_injection.PostCookiesCommandHandlerFunc(func(params cmd_injection.PostCookiesCommandParams) middleware.Responder {
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			var data string

			cookie := http.Cookie{Name: "cookie", Value:"value"}
			http.SetCookie(w, &cookie)

			txt, render := cmdi.ExecHandlerCtx(params.Safety, "unusable")
			if !render {
				data = string(txt)
			} else {
				data = "not intended to render"
			}

			if err := p.Produce(w, data); err != nil {
			}
		})

	})
	api.CmdInjectionPostCookiesCommandContextHandler = cmd_injection.PostCookiesCommandContextHandlerFunc(func(params cmd_injection.PostCookiesCommandContextParams) middleware.Responder {
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			var data string

			cookie := http.Cookie{Name: "cookie", Value:"value"}
			http.SetCookie(w, &cookie)

			txt, render := cmdi.ExecHandlerCtx(params.Safety, "unusable")
			if !render {
				data = string(txt)
			} else {
				data = "not intended to render"
			}

			if err := p.Produce(w, data); err != nil {
			}
		})

	})

	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Go Swagger API integrated with Go Test Bench"
	parser.LongDescription = "An API built with go-swagger to generate intentionally vulnerable endpoints"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}

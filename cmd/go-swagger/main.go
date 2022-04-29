package main

import (
	"log"
	"os"

	// only if we are going to be implementing runtime.Producer.
	//"github.com/go-openapi/runtime"

	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"

	"github.com/go-openapi/runtime/middleware"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"


	"github.com/go-openapi/loads"

	flags "github.com/jessevdk/go-flags"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
)

func main() {

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewSwaggerBenchAPI(swaggerSpec)

	//api.HTMLProducer = runtime.TextProducer()

	api.SwaggerServerRootHandler = swagger_server.RootHandlerFunc(serveswagger.SwaggerRootHandler)

	api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(serveswagger.CommandInjectionHandler)

	api.CmdInjectionGetQueryExploitHandler = cmd_injection.GetQueryExploitHandlerFunc(func(params cmd_injection.GetQueryExploitParams) middleware.Responder {
		var payload string

		if params.Command == "exec.Command" {
			txt, render := cmdi.ExecHandler(params.Safety, params.Input)
			if !render {
				payload = string(txt)
			} else {
				payload = "exec.Command operation failure"
			}
			return cmd_injection.NewGetQueryExploitOK().WithPayload(payload)
		} else if params.Command == "exec.CommandContext" {
			txt, render := cmdi.ExecHandlerCtx(params.Safety, params.Input)
			if !render {
				payload = string(txt)
			} else {
				payload = "exec.CommandContext operation failure"
			}
		}
		return cmd_injection.NewGetQueryExploitOK().WithPayload(payload)
	})

	api.CmdInjectionPostCookiesExploitHandler = cmd_injection.PostCookiesExploitHandlerFunc(func(params cmd_injection.PostCookiesExploitParams) middleware.Responder {
		//TODO: Something still needs to be resolved:
		// 1) Setting the yml to allow passing the actual cookie as parameter

		var payload string

		if params.Command == "exec.Command" {
			txt, render := cmdi.ExecHandler(params.Safety, params.Input)
			if !render {
				payload = string(txt)
			} else {
				payload = "exec.Command operation failure"
			}
			return cmd_injection.NewGetQueryExploitOK().WithPayload(payload)
		} else if params.Command == "exec.CommandContext" {
			txt, render := cmdi.ExecHandlerCtx(params.Safety, params.Input)
			if !render {
				payload = string(txt)
			} else {
				payload = "exec.CommandContext operation failure"
			}
		}
		return cmd_injection.NewPostCookiesExploitOK().WithPayload(payload)
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

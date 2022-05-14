// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/path_traversal"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/sql_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/unvalidated_redirect"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/xss"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"
)

//go:generate swagger generate server --target ../../go-swagger --name SwaggerBench --spec ../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.SwaggerBenchAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SwaggerBenchAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.HTMLProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		var (
			t                  = common.Templates["underConstruction.gohtml"]
			params interface{} = serveswagger.SwaggerParams
		)
		if str, ok := data.(string); ok {
			for _, r := range common.AllRoutes {
				log.Println("loading template file:", r.TmplFile)
				log.Println("route Base:", r.Base)
				if str != r.Base {
					continue
				}
				tmpl, ok := common.Templates[r.TmplFile]
				if !ok {
					break
				}
				t = tmpl
				params = common.Parameters{
					ConstParams: serveswagger.SwaggerParams,
					Name:        r.Base,
				}
			}
		}

		t.ExecuteTemplate(w, "layout.gohtml", params)
		return nil
	})

	api.JSONProducer = runtime.JSONProducer()
	api.TxtProducer = runtime.TextProducer()

	if api.CmdInjectionCmdInjectionFrontHandler == nil {
		api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(func(params cmd_injection.CmdInjectionFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.CmdInjectionFront has not yet been implemented")
		})
	}
	if api.CmdInjectionGetQueryCommandHandler == nil {
		api.CmdInjectionGetQueryCommandHandler = cmd_injection.GetQueryCommandHandlerFunc(func(params cmd_injection.GetQueryCommandParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.GetQueryCommand has not yet been implemented")
		})
	}
	if api.CmdInjectionGetQueryCommandContextHandler == nil {
		api.CmdInjectionGetQueryCommandContextHandler = cmd_injection.GetQueryCommandContextHandlerFunc(func(params cmd_injection.GetQueryCommandContextParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.GetQueryCommandContext has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalFrontHandler == nil {
		api.PathTraversalPathTraversalFrontHandler = path_traversal.PathTraversalFrontHandlerFunc(func(params path_traversal.PathTraversalFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalFront has not yet been implemented")
		})
	}
	if api.SwaggerServerRootHandler == nil {
		api.SwaggerServerRootHandler = swagger_server.RootHandlerFunc(func(params swagger_server.RootParams) middleware.Responder {
			return middleware.NotImplemented("operation swagger_server.Root has not yet been implemented")
		})
	}
	if api.SQLInjectionSQLInjectionFrontHandler == nil {
		api.SQLInjectionSQLInjectionFrontHandler = sql_injection.SQLInjectionFrontHandlerFunc(func(params sql_injection.SQLInjectionFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation sql_injection.SQLInjectionFront has not yet been implemented")
		})
	}
	if api.SsrfSsrfFrontHandler == nil {
		api.SsrfSsrfFrontHandler = ssrf.SsrfFrontHandlerFunc(func(params ssrf.SsrfFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation ssrf.SsrfFront has not yet been implemented")
		})
	}
	if api.UnvalidatedRedirectUnvalidatedRedirectFrontHandler == nil {
		api.UnvalidatedRedirectUnvalidatedRedirectFrontHandler = unvalidated_redirect.UnvalidatedRedirectFrontHandlerFunc(func(params unvalidated_redirect.UnvalidatedRedirectFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation unvalidated_redirect.UnvalidatedRedirectFront has not yet been implemented")
		})
	}
	if api.XSSXSSFrontHandler == nil {
		api.XSSXSSFrontHandler = xss.XSSFrontHandlerFunc(func(params xss.XSSFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation xss.XSSFront has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

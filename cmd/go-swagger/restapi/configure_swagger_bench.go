// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/path_traversal"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/sql_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/xss"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

//go:generate go run ../regen/regen.go
//go:generate rm -rf ./operations/

//original generate command: //go:generate swagger generate server --target ../../go-swagger --name SwaggerBench --spec ../swagger.yml --principal interface{} --exclude-main

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

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	if api.CmdInjectionCmdInjectionFrontHandler == nil {
		api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(func(params cmd_injection.CmdInjectionFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.CmdInjectionFront has not yet been implemented")
		})
	}
	if api.CmdInjectionCmdInjectionGetQueryCommandHandler == nil {
		api.CmdInjectionCmdInjectionGetQueryCommandHandler = cmd_injection.CmdInjectionGetQueryCommandHandlerFunc(func(params cmd_injection.CmdInjectionGetQueryCommandParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.CmdInjectionGetQueryCommand has not yet been implemented")
		})
	}
	if api.CmdInjectionCmdInjectionGetQueryCommandContextHandler == nil {
		api.CmdInjectionCmdInjectionGetQueryCommandContextHandler = cmd_injection.CmdInjectionGetQueryCommandContextHandlerFunc(func(params cmd_injection.CmdInjectionGetQueryCommandContextParams) middleware.Responder {
			return middleware.NotImplemented("operation cmd_injection.CmdInjectionGetQueryCommandContext has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalFrontHandler == nil {
		api.PathTraversalPathTraversalFrontHandler = path_traversal.PathTraversalFrontHandlerFunc(func(params path_traversal.PathTraversalFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalFront has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetBufferedQueryCreateHandler == nil {
		api.PathTraversalPathTraversalGetBufferedQueryCreateHandler = path_traversal.PathTraversalGetBufferedQueryCreateHandlerFunc(func(params path_traversal.PathTraversalGetBufferedQueryCreateParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetBufferedQueryCreate has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetBufferedQueryOpenHandler == nil {
		api.PathTraversalPathTraversalGetBufferedQueryOpenHandler = path_traversal.PathTraversalGetBufferedQueryOpenHandlerFunc(func(params path_traversal.PathTraversalGetBufferedQueryOpenParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetBufferedQueryOpen has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetBufferedQueryReadFileHandler == nil {
		api.PathTraversalPathTraversalGetBufferedQueryReadFileHandler = path_traversal.PathTraversalGetBufferedQueryReadFileHandlerFunc(func(params path_traversal.PathTraversalGetBufferedQueryReadFileParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetBufferedQueryReadFile has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetBufferedQueryWriteFileHandler == nil {
		api.PathTraversalPathTraversalGetBufferedQueryWriteFileHandler = path_traversal.PathTraversalGetBufferedQueryWriteFileHandlerFunc(func(params path_traversal.PathTraversalGetBufferedQueryWriteFileParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetBufferedQueryWriteFile has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetQueryCreateHandler == nil {
		api.PathTraversalPathTraversalGetQueryCreateHandler = path_traversal.PathTraversalGetQueryCreateHandlerFunc(func(params path_traversal.PathTraversalGetQueryCreateParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetQueryCreate has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetQueryOpenHandler == nil {
		api.PathTraversalPathTraversalGetQueryOpenHandler = path_traversal.PathTraversalGetQueryOpenHandlerFunc(func(params path_traversal.PathTraversalGetQueryOpenParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetQueryOpen has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetQueryReadFileHandler == nil {
		api.PathTraversalPathTraversalGetQueryReadFileHandler = path_traversal.PathTraversalGetQueryReadFileHandlerFunc(func(params path_traversal.PathTraversalGetQueryReadFileParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetQueryReadFile has not yet been implemented")
		})
	}
	if api.PathTraversalPathTraversalGetQueryWriteFileHandler == nil {
		api.PathTraversalPathTraversalGetQueryWriteFileHandler = path_traversal.PathTraversalGetQueryWriteFileHandlerFunc(func(params path_traversal.PathTraversalGetQueryWriteFileParams) middleware.Responder {
			return middleware.NotImplemented("operation path_traversal.PathTraversalGetQueryWriteFile has not yet been implemented")
		})
	}
	if api.SQLInjectionSQLInjectionFrontHandler == nil {
		api.SQLInjectionSQLInjectionFrontHandler = sql_injection.SQLInjectionFrontHandlerFunc(func(params sql_injection.SQLInjectionFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation sql_injection.SQLInjectionFront has not yet been implemented")
		})
	}
	if api.SQLInjectionSQLInjectionGetQueryExecHandler == nil {
		api.SQLInjectionSQLInjectionGetQueryExecHandler = sql_injection.SQLInjectionGetQueryExecHandlerFunc(func(params sql_injection.SQLInjectionGetQueryExecParams) middleware.Responder {
			return middleware.NotImplemented("operation sql_injection.SQLInjectionGetQueryExec has not yet been implemented")
		})
	}
	if api.SsrfSsrfFrontHandler == nil {
		api.SsrfSsrfFrontHandler = ssrf.SsrfFrontHandlerFunc(func(params ssrf.SsrfFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation ssrf.SsrfFront has not yet been implemented")
		})
	}
	if api.SsrfSsrfGetQueryHTTPHandler == nil {
		api.SsrfSsrfGetQueryHTTPHandler = ssrf.SsrfGetQueryHTTPHandlerFunc(func(params ssrf.SsrfGetQueryHTTPParams) middleware.Responder {
			return middleware.NotImplemented("operation ssrf.SsrfGetQuerySink has not yet been implemented")
		})
	}
	if api.XSSXSSFrontHandler == nil {
		api.XSSXSSFrontHandler = xss.XSSFrontHandlerFunc(func(params xss.XSSFrontParams) middleware.Responder {
			return middleware.NotImplemented("operation xss.XSSFront has not yet been implemented")
		})
	}
	if api.XSSXSSGetBufferedQueryReflectedXSSHandler == nil {
		api.XSSXSSGetBufferedQueryReflectedXSSHandler = xss.XSSGetBufferedQueryReflectedXSSHandlerFunc(func(params xss.XSSGetBufferedQueryReflectedXSSParams) middleware.Responder {
			return middleware.NotImplemented("operation xss.XSSGetBufferedQuerySink has not yet been implemented")
		})
	}
	if api.XSSXSSGetQueryReflectedXSSHandler == nil {
		api.XSSXSSGetQueryReflectedXSSHandler = xss.XSSGetQueryReflectedXSSHandlerFunc(func(params xss.XSSGetQueryReflectedXSSParams) middleware.Responder {
			return middleware.NotImplemented("operation xss.XSSGetQuerySink has not yet been implemented")
		})
	}
	if api.SwaggerServerRootHandler == nil {
		api.SwaggerServerRootHandler = swagger_server.RootHandlerFunc(func(params swagger_server.RootParams) middleware.Responder {
			return middleware.NotImplemented("operation swagger_server.Root has not yet been implemented")
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

var assetHandler http.Handler = http.StripPrefix("/assets/",
	http.FileServer(http.Dir(
		func() string {
			dir, err := common.LocateDir("public", 5)
			if err != nil {
				log.Fatal(err)
			}
			return dir
		}())),
)

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			assetHandler.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return uiMiddleware(handler)
}

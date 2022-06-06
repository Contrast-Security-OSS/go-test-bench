// GENERATED CODE - DO NOT EDIT
// Generated at 2022-06-06T19:32:24Z. To re-generate, run the following in the repo root:
// go run ./cmd/go-swagger/regen/regen.go

package serveswagger

import (
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/path_traversal"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/sql_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/unvalidated_redirect"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/xss"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/go-openapi/runtime/middleware"
)

func generatedInit(api *operations.SwaggerBenchAPI, rmap common.RouteMap, pd *common.ConstParams) {

	// CmdInjection

	api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(
		func(p cmd_injection.CmdInjectionFrontParams) middleware.Responder {
			return RouteHandler(rmap["/cmdInjection"], pd, p.HTTPRequest)
		},
	)
	api.CmdInjectionCmdInjectionGetQueryCommandHandler = cmd_injection.CmdInjectionGetQueryCommandHandlerFunc(
		func(p cmd_injection.CmdInjectionGetQueryCommandParams) middleware.Responder {
			return RouteHandler(rmap["/cmdInjection"], pd, p.HTTPRequest)
		},
	)
	api.CmdInjectionCmdInjectionGetQueryCommandContextHandler = cmd_injection.CmdInjectionGetQueryCommandContextHandlerFunc(
		func(p cmd_injection.CmdInjectionGetQueryCommandContextParams) middleware.Responder {
			return RouteHandler(rmap["/cmdInjection"], pd, p.HTTPRequest)
		},
	)

	// PathTraversal

	api.PathTraversalPathTraversalFrontHandler = path_traversal.PathTraversalFrontHandlerFunc(
		func(p path_traversal.PathTraversalFrontParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetQueryReadFileHandler = path_traversal.PathTraversalGetQueryReadFileHandlerFunc(
		func(p path_traversal.PathTraversalGetQueryReadFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetQueryOpenHandler = path_traversal.PathTraversalGetQueryOpenHandlerFunc(
		func(p path_traversal.PathTraversalGetQueryOpenParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetQueryWriteFileHandler = path_traversal.PathTraversalGetQueryWriteFileHandlerFunc(
		func(p path_traversal.PathTraversalGetQueryWriteFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetQueryCreateHandler = path_traversal.PathTraversalGetQueryCreateHandlerFunc(
		func(p path_traversal.PathTraversalGetQueryCreateParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBufferedQueryReadFileHandler = path_traversal.PathTraversalGetBufferedQueryReadFileHandlerFunc(
		func(p path_traversal.PathTraversalGetBufferedQueryReadFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBufferedQueryOpenHandler = path_traversal.PathTraversalGetBufferedQueryOpenHandlerFunc(
		func(p path_traversal.PathTraversalGetBufferedQueryOpenParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBufferedQueryWriteFileHandler = path_traversal.PathTraversalGetBufferedQueryWriteFileHandlerFunc(
		func(p path_traversal.PathTraversalGetBufferedQueryWriteFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBufferedQueryCreateHandler = path_traversal.PathTraversalGetBufferedQueryCreateHandlerFunc(
		func(p path_traversal.PathTraversalGetBufferedQueryCreateParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)

	// XSS

	api.XSSXSSFrontHandler = xss.XSSFrontHandlerFunc(
		func(p xss.XSSFrontParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)
	api.XSSXSSGetQuerySinkHandler = xss.XSSGetQuerySinkHandlerFunc(
		func(p xss.XSSGetQuerySinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)
	api.XSSXSSGetBufferedQuerySinkHandler = xss.XSSGetBufferedQuerySinkHandlerFunc(
		func(p xss.XSSGetBufferedQuerySinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)

	// SQLInjection

	api.SQLInjectionSQLInjectionFrontHandler = sql_injection.SQLInjectionFrontHandlerFunc(
		func(p sql_injection.SQLInjectionFrontParams) middleware.Responder {
			return RouteHandler(rmap["/sqlInjection"], pd, p.HTTPRequest)
		},
	)
	api.SQLInjectionSQLInjectionGetQueryExecHandler = sql_injection.SQLInjectionGetQueryExecHandlerFunc(
		func(p sql_injection.SQLInjectionGetQueryExecParams) middleware.Responder {
			return RouteHandler(rmap["/sqlInjection"], pd, p.HTTPRequest)
		},
	)

	// Ssrf

	api.SsrfSsrfFrontHandler = ssrf.SsrfFrontHandlerFunc(
		func(p ssrf.SsrfFrontParams) middleware.Responder {
			return RouteHandler(rmap["/ssrf"], pd, p.HTTPRequest)
		},
	)
	api.SsrfSsrfGetQuerySinkHandler = ssrf.SsrfGetQuerySinkHandlerFunc(
		func(p ssrf.SsrfGetQuerySinkParams) middleware.Responder {
			return RouteHandler(rmap["/ssrf"], pd, p.HTTPRequest)
		},
	)

	// UnvalidatedRedirect

	api.UnvalidatedRedirectUnvalidatedRedirectFrontHandler = unvalidated_redirect.UnvalidatedRedirectFrontHandlerFunc(
		func(p unvalidated_redirect.UnvalidatedRedirectFrontParams) middleware.Responder {
			return RouteHandler(rmap["/unvalidatedRedirect"], pd, p.HTTPRequest)
		},
	)
	api.UnvalidatedRedirectUnvalidatedRedirectGetQueryRedirectHandler = unvalidated_redirect.UnvalidatedRedirectGetQueryRedirectHandlerFunc(
		func(p unvalidated_redirect.UnvalidatedRedirectGetQueryRedirectParams) middleware.Responder {
			return RouteHandler(rmap["/unvalidatedRedirect"], pd, p.HTTPRequest)
		},
	)
}

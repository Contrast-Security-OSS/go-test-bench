// GENERATED CODE - DO NOT EDIT
// Generated at 2022-05-24T15:44:12Z. To re-generate, run the following in the repo root:
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

func generatedInit(api *operations.SwaggerBenchAPI, rmap common.RouteMap, pd common.ConstParams) {

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
	api.CmdInjectionCmdInjectionGetCookiesCommandHandler = cmd_injection.CmdInjectionGetCookiesCommandHandlerFunc(
		func(p cmd_injection.CmdInjectionGetCookiesCommandParams) middleware.Responder {
			return RouteHandler(rmap["/cmdInjection"], pd, p.HTTPRequest)
		},
	)
	api.CmdInjectionCmdInjectionGetCookiesCommandContextHandler = cmd_injection.CmdInjectionGetCookiesCommandContextHandlerFunc(
		func(p cmd_injection.CmdInjectionGetCookiesCommandContextParams) middleware.Responder {
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
	api.PathTraversalPathTraversalGetHeadersReadFileHandler = path_traversal.PathTraversalGetHeadersReadFileHandlerFunc(
		func(p path_traversal.PathTraversalGetHeadersReadFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetHeadersOpenHandler = path_traversal.PathTraversalGetHeadersOpenHandlerFunc(
		func(p path_traversal.PathTraversalGetHeadersOpenParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetHeadersWriteFileHandler = path_traversal.PathTraversalGetHeadersWriteFileHandlerFunc(
		func(p path_traversal.PathTraversalGetHeadersWriteFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetHeadersCreateHandler = path_traversal.PathTraversalGetHeadersCreateHandlerFunc(
		func(p path_traversal.PathTraversalGetHeadersCreateParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBodyReadFileHandler = path_traversal.PathTraversalGetBodyReadFileHandlerFunc(
		func(p path_traversal.PathTraversalGetBodyReadFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBodyOpenHandler = path_traversal.PathTraversalGetBodyOpenHandlerFunc(
		func(p path_traversal.PathTraversalGetBodyOpenParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBodyWriteFileHandler = path_traversal.PathTraversalGetBodyWriteFileHandlerFunc(
		func(p path_traversal.PathTraversalGetBodyWriteFileParams) middleware.Responder {
			return RouteHandler(rmap["/pathTraversal"], pd, p.HTTPRequest)
		},
	)
	api.PathTraversalPathTraversalGetBodyCreateHandler = path_traversal.PathTraversalGetBodyCreateHandlerFunc(
		func(p path_traversal.PathTraversalGetBodyCreateParams) middleware.Responder {
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
	api.XSSXSSGetParamsSinkHandler = xss.XSSGetParamsSinkHandlerFunc(
		func(p xss.XSSGetParamsSinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)
	api.XSSXSSGetBodySinkHandler = xss.XSSGetBodySinkHandlerFunc(
		func(p xss.XSSGetBodySinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)
	api.XSSXSSGetBufferedBodySinkHandler = xss.XSSGetBufferedBodySinkHandlerFunc(
		func(p xss.XSSGetBufferedBodySinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)
	api.XSSXSSGetResponseSinkHandler = xss.XSSGetResponseSinkHandlerFunc(
		func(p xss.XSSGetResponseSinkParams) middleware.Responder {
			return RouteHandler(rmap["/xss"], pd, p.HTTPRequest)
		},
	)

	// SQLInjection

	api.SQLInjectionSQLInjectionFrontHandler = sql_injection.SQLInjectionFrontHandlerFunc(
		func(p sql_injection.SQLInjectionFrontParams) middleware.Responder {
			return RouteHandler(rmap["/sqlInjection"], pd, p.HTTPRequest)
		},
	)
	api.SQLInjectionSQLInjectionGetBodyExecHandler = sql_injection.SQLInjectionGetBodyExecHandlerFunc(
		func(p sql_injection.SQLInjectionGetBodyExecParams) middleware.Responder {
			return RouteHandler(rmap["/sqlInjection"], pd, p.HTTPRequest)
		},
	)
	api.SQLInjectionSQLInjectionGetQueryExecHandler = sql_injection.SQLInjectionGetQueryExecHandlerFunc(
		func(p sql_injection.SQLInjectionGetQueryExecParams) middleware.Responder {
			return RouteHandler(rmap["/sqlInjection"], pd, p.HTTPRequest)
		},
	)
	api.SQLInjectionSQLInjectionGetHeadersJSONExecHandler = sql_injection.SQLInjectionGetHeadersJSONExecHandlerFunc(
		func(p sql_injection.SQLInjectionGetHeadersJSONExecParams) middleware.Responder {
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
	api.SsrfSsrfGetPathSinkHandler = ssrf.SsrfGetPathSinkHandlerFunc(
		func(p ssrf.SsrfGetPathSinkParams) middleware.Responder {
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

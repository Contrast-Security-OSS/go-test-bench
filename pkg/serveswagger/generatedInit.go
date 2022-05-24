// GENERATED CODE - DO NOT EDIT
// Generated at 2022-05-24T15:26:21Z. To re-generate, run the following in the repo root:
// go run ./cmd/go-swagger/regen/regen.go

package serveswagger

import (
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
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

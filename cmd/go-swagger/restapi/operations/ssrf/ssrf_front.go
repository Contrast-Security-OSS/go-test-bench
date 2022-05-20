// Code generated by go-swagger; DO NOT EDIT.

package ssrf

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SsrfFrontHandlerFunc turns a function with the right signature into a ssrf front handler
type SsrfFrontHandlerFunc func(SsrfFrontParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SsrfFrontHandlerFunc) Handle(params SsrfFrontParams) middleware.Responder {
	return fn(params)
}

// SsrfFrontHandler interface for that can handle valid ssrf front params
type SsrfFrontHandler interface {
	Handle(SsrfFrontParams) middleware.Responder
}

// NewSsrfFront creates a new http.Handler for the ssrf front operation
func NewSsrfFront(ctx *middleware.Context, handler SsrfFrontHandler) *SsrfFront {
	return &SsrfFront{Context: ctx, Handler: handler}
}

/* SsrfFront swagger:route GET /ssrf ssrf ssrfFront

front page of the SSRF vulnerability

*/
type SsrfFront struct {
	Context *middleware.Context
	Handler SsrfFrontHandler
}

func (o *SsrfFront) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSsrfFrontParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

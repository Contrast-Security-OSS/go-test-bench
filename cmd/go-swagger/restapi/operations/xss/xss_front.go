// Code generated by go-swagger; DO NOT EDIT.

package xss

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// XSSFrontHandlerFunc turns a function with the right signature into a xss front handler
type XSSFrontHandlerFunc func(XSSFrontParams) middleware.Responder

// Handle executing the request and returning a response
func (fn XSSFrontHandlerFunc) Handle(params XSSFrontParams) middleware.Responder {
	return fn(params)
}

// XSSFrontHandler interface for that can handle valid xss front params
type XSSFrontHandler interface {
	Handle(XSSFrontParams) middleware.Responder
}

// NewXSSFront creates a new http.Handler for the xss front operation
func NewXSSFront(ctx *middleware.Context, handler XSSFrontHandler) *XSSFront {
	return &XSSFront{Context: ctx, Handler: handler}
}

/* XSSFront swagger:route GET /xss xss xssFront

supposed to serve the frontend for the query or cookie vulns

*/
type XSSFront struct {
	Context *middleware.Context
	Handler XSSFrontHandler
}

func (o *XSSFront) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewXSSFrontParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

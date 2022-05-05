// Code generated by go-swagger; DO NOT EDIT.

package path_traversal

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PathTraversaslFrontOKCode is the HTTP code returned for type PathTraversaslFrontOK
const PathTraversaslFrontOKCode int = 200

/*PathTraversaslFrontOK served front end for path traversal page of Swagger API

swagger:response pathTraversaslFrontOK
*/
type PathTraversaslFrontOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPathTraversaslFrontOK creates PathTraversaslFrontOK with default headers values
func NewPathTraversaslFrontOK() *PathTraversaslFrontOK {

	return &PathTraversaslFrontOK{}
}

// WithPayload adds the payload to the path traversasl front o k response
func (o *PathTraversaslFrontOK) WithPayload(payload string) *PathTraversaslFrontOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the path traversasl front o k response
func (o *PathTraversaslFrontOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PathTraversaslFrontOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*PathTraversaslFrontDefault error occured

swagger:response pathTraversaslFrontDefault
*/
type PathTraversaslFrontDefault struct {
	_statusCode int
}

// NewPathTraversaslFrontDefault creates PathTraversaslFrontDefault with default headers values
func NewPathTraversaslFrontDefault(code int) *PathTraversaslFrontDefault {
	if code <= 0 {
		code = 500
	}

	return &PathTraversaslFrontDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the path traversasl front default response
func (o *PathTraversaslFrontDefault) WithStatusCode(code int) *PathTraversaslFrontDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the path traversasl front default response
func (o *PathTraversaslFrontDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WriteResponse to the client
func (o *PathTraversaslFrontDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(o._statusCode)
}
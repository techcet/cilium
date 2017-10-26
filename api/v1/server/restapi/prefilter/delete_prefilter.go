// Code generated by go-swagger; DO NOT EDIT.

package prefilter

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeletePrefilterHandlerFunc turns a function with the right signature into a delete prefilter handler
type DeletePrefilterHandlerFunc func(DeletePrefilterParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeletePrefilterHandlerFunc) Handle(params DeletePrefilterParams) middleware.Responder {
	return fn(params)
}

// DeletePrefilterHandler interface for that can handle valid delete prefilter params
type DeletePrefilterHandler interface {
	Handle(DeletePrefilterParams) middleware.Responder
}

// NewDeletePrefilter creates a new http.Handler for the delete prefilter operation
func NewDeletePrefilter(ctx *middleware.Context, handler DeletePrefilterHandler) *DeletePrefilter {
	return &DeletePrefilter{Context: ctx, Handler: handler}
}

/*DeletePrefilter swagger:route DELETE /prefilter prefilter deletePrefilter

Delete list of CIDRs

*/
type DeletePrefilter struct {
	Context *middleware.Context
	Handler DeletePrefilterHandler
}

func (o *DeletePrefilter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeletePrefilterParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

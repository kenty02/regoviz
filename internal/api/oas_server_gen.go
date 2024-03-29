// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// AstGet implements GET /ast operation.
	//
	// GET /ast
	AstGet(ctx context.Context, params AstGetParams) (*AstGetOK, error)
	// AstPrettyGet implements GET /astPretty operation.
	//
	// GET /astPretty
	AstPrettyGet(ctx context.Context, params AstPrettyGetParams) (*AstPrettyGetOK, error)
	// CallTreeAvailableEntrypointsGet implements GET /callTree/availableEntrypoints operation.
	//
	// GET /callTree/availableEntrypoints
	CallTreeAvailableEntrypointsGet(ctx context.Context, params CallTreeAvailableEntrypointsGetParams) (*CallTreeAvailableEntrypointsGetOK, error)
	// CallTreeGet implements GET /callTree operation.
	//
	// GET /callTree
	CallTreeGet(ctx context.Context, params CallTreeGetParams) (*CallTreeGetOK, error)
	// DepTreeTextGet implements GET /depTreeText operation.
	//
	// GET /depTreeText
	DepTreeTextGet(ctx context.Context, params DepTreeTextGetParams) (*DepTreeTextGetOK, error)
	// FlowchartGet implements GET /flowchart operation.
	//
	// GET /flowchart
	FlowchartGet(ctx context.Context, params FlowchartGetParams) (*FlowchartGetOK, error)
	// IrGet implements GET /ir operation.
	//
	// GET /ir
	IrGet(ctx context.Context, params IrGetParams) (*IrGetOK, error)
	// SamplesGet implements GET /samples operation.
	//
	// GET /samples
	SamplesGet(ctx context.Context) ([]Sample, error)
	// VarTracePost implements POST /varTrace operation.
	//
	// POST /varTrace
	VarTracePost(ctx context.Context, params VarTracePostParams) (*VarTracePostOK, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}

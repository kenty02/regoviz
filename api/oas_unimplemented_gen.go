// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AstGet implements GET /ast operation.
//
// GET /ast
func (UnimplementedHandler) AstGet(ctx context.Context, params AstGetParams) (r *AstGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// CallTreeGet implements GET /callTree operation.
//
// GET /callTree
func (UnimplementedHandler) CallTreeGet(ctx context.Context, params CallTreeGetParams) (r *CallTreeGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// DepTreeTextGet implements GET /depTreeText operation.
//
// GET /depTreeText
func (UnimplementedHandler) DepTreeTextGet(ctx context.Context, params DepTreeTextGetParams) (r *DepTreeTextGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// FlowchartGet implements GET /flowchart operation.
//
// GET /flowchart
func (UnimplementedHandler) FlowchartGet(ctx context.Context, params FlowchartGetParams) (r *FlowchartGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// IrGet implements GET /ir operation.
//
// GET /ir
func (UnimplementedHandler) IrGet(ctx context.Context, params IrGetParams) (r *IrGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// RulesGet implements GET /rules operation.
//
// GET /rules
func (UnimplementedHandler) RulesGet(ctx context.Context) (r []Rule, _ error) {
	return r, ht.ErrNotImplemented
}

// SamplesGet implements GET /samples operation.
//
// GET /samples
func (UnimplementedHandler) SamplesGet(ctx context.Context) (r []Sample, _ error) {
	return r, ht.ErrNotImplemented
}

// VarTracePost implements POST /varTrace operation.
//
// POST /varTrace
func (UnimplementedHandler) VarTracePost(ctx context.Context, params VarTracePostParams) (r *VarTracePostOK, _ error) {
	return r, ht.ErrNotImplemented
}
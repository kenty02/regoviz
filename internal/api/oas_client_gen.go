// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
)

// Invoker invokes operations described by OpenAPI v3 specification.
type Invoker interface {
	// AstGet invokes GET /ast operation.
	//
	// GET /ast
	AstGet(ctx context.Context, params AstGetParams) (*AstGetOK, error)
	// AstPrettyGet invokes GET /astPretty operation.
	//
	// GET /astPretty
	AstPrettyGet(ctx context.Context, params AstPrettyGetParams) (*AstPrettyGetOK, error)
	// CallTreeAvailableEntrypointsGet invokes GET /callTree/availableEntrypoints operation.
	//
	// GET /callTree/availableEntrypoints
	CallTreeAvailableEntrypointsGet(ctx context.Context, params CallTreeAvailableEntrypointsGetParams) (*CallTreeAvailableEntrypointsGetOK, error)
	// CallTreeGet invokes GET /callTree operation.
	//
	// GET /callTree
	CallTreeGet(ctx context.Context, params CallTreeGetParams) (*CallTreeGetOK, error)
	// DepTreeTextGet invokes GET /depTreeText operation.
	//
	// GET /depTreeText
	DepTreeTextGet(ctx context.Context, params DepTreeTextGetParams) (*DepTreeTextGetOK, error)
	// FlowchartGet invokes GET /flowchart operation.
	//
	// GET /flowchart
	FlowchartGet(ctx context.Context, params FlowchartGetParams) (*FlowchartGetOK, error)
	// IrGet invokes GET /ir operation.
	//
	// GET /ir
	IrGet(ctx context.Context, params IrGetParams) (*IrGetOK, error)
	// SamplesGet invokes GET /samples operation.
	//
	// GET /samples
	SamplesGet(ctx context.Context) ([]Sample, error)
	// VarTracePost invokes POST /varTrace operation.
	//
	// POST /varTrace
	VarTracePost(ctx context.Context, params VarTracePostParams) (*VarTracePostOK, error)
}

// Client implements OAS client.
type Client struct {
	serverURL *url.URL
	sec       SecuritySource
	baseClient
}

var _ Handler = struct {
	*Client
}{}

func trimTrailingSlashes(u *url.URL) {
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
}

// NewClient initializes new Client defined by OAS.
func NewClient(serverURL string, sec SecuritySource, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	trimTrailingSlashes(u)

	c, err := newClientConfig(opts...).baseClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		serverURL:  u,
		sec:        sec,
		baseClient: c,
	}, nil
}

type serverURLKey struct{}

// WithServerURL sets context key to override server URL.
func WithServerURL(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, serverURLKey{}, u)
}

func (c *Client) requestURL(ctx context.Context) *url.URL {
	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
	if !ok {
		return c.serverURL
	}
	return u
}

// AstGet invokes GET /ast operation.
//
// GET /ast
func (c *Client) AstGet(ctx context.Context, params AstGetParams) (*AstGetOK, error) {
	res, err := c.sendAstGet(ctx, params)
	return res, err
}

func (c *Client) sendAstGet(ctx context.Context, params AstGetParams) (res *AstGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/ast"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "AstGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/ast"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "AstGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeAstGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AstPrettyGet invokes GET /astPretty operation.
//
// GET /astPretty
func (c *Client) AstPrettyGet(ctx context.Context, params AstPrettyGetParams) (*AstPrettyGetOK, error) {
	res, err := c.sendAstPrettyGet(ctx, params)
	return res, err
}

func (c *Client) sendAstPrettyGet(ctx context.Context, params AstPrettyGetParams) (res *AstPrettyGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/astPretty"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "AstPrettyGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/astPretty"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "AstPrettyGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeAstPrettyGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// CallTreeAvailableEntrypointsGet invokes GET /callTree/availableEntrypoints operation.
//
// GET /callTree/availableEntrypoints
func (c *Client) CallTreeAvailableEntrypointsGet(ctx context.Context, params CallTreeAvailableEntrypointsGetParams) (*CallTreeAvailableEntrypointsGetOK, error) {
	res, err := c.sendCallTreeAvailableEntrypointsGet(ctx, params)
	return res, err
}

func (c *Client) sendCallTreeAvailableEntrypointsGet(ctx context.Context, params CallTreeAvailableEntrypointsGetParams) (res *CallTreeAvailableEntrypointsGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/callTree/availableEntrypoints"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "CallTreeAvailableEntrypointsGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/callTree/availableEntrypoints"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "CallTreeAvailableEntrypointsGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeCallTreeAvailableEntrypointsGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// CallTreeGet invokes GET /callTree operation.
//
// GET /callTree
func (c *Client) CallTreeGet(ctx context.Context, params CallTreeGetParams) (*CallTreeGetOK, error) {
	res, err := c.sendCallTreeGet(ctx, params)
	return res, err
}

func (c *Client) sendCallTreeGet(ctx context.Context, params CallTreeGetParams) (res *CallTreeGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/callTree"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "CallTreeGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/callTree"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "entrypoint" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "entrypoint",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Entrypoint))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "input" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "input",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Input.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "data" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "data",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Data.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "query" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "query",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Query.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "CallTreeGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeCallTreeGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// DepTreeTextGet invokes GET /depTreeText operation.
//
// GET /depTreeText
func (c *Client) DepTreeTextGet(ctx context.Context, params DepTreeTextGetParams) (*DepTreeTextGetOK, error) {
	res, err := c.sendDepTreeTextGet(ctx, params)
	return res, err
}

func (c *Client) sendDepTreeTextGet(ctx context.Context, params DepTreeTextGetParams) (res *DepTreeTextGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/depTreeText"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "DepTreeTextGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/depTreeText"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "DepTreeTextGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeDepTreeTextGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// FlowchartGet invokes GET /flowchart operation.
//
// GET /flowchart
func (c *Client) FlowchartGet(ctx context.Context, params FlowchartGetParams) (*FlowchartGetOK, error) {
	res, err := c.sendFlowchartGet(ctx, params)
	return res, err
}

func (c *Client) sendFlowchartGet(ctx context.Context, params FlowchartGetParams) (res *FlowchartGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/flowchart"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "FlowchartGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/flowchart"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "edit" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "edit",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Edit.Get(); ok {
				return e.EncodeValue(conv.BoolToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "FlowchartGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeFlowchartGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// IrGet invokes GET /ir operation.
//
// GET /ir
func (c *Client) IrGet(ctx context.Context, params IrGetParams) (*IrGetOK, error) {
	res, err := c.sendIrGet(ctx, params)
	return res, err
}

func (c *Client) sendIrGet(ctx context.Context, params IrGetParams) (res *IrGetOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/ir"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "IrGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/ir"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "IrGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeIrGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// SamplesGet invokes GET /samples operation.
//
// GET /samples
func (c *Client) SamplesGet(ctx context.Context) ([]Sample, error) {
	res, err := c.sendSamplesGet(ctx)
	return res, err
}

func (c *Client) sendSamplesGet(ctx context.Context) (res []Sample, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/samples"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "SamplesGet",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/samples"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "SamplesGet", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeSamplesGetResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// VarTracePost invokes POST /varTrace operation.
//
// POST /varTrace
func (c *Client) VarTracePost(ctx context.Context, params VarTracePostParams) (*VarTracePostOK, error) {
	res, err := c.sendVarTracePost(ctx, params)
	return res, err
}

func (c *Client) sendVarTracePost(ctx context.Context, params VarTracePostParams) (res *VarTracePostOK, err error) {
	otelAttrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/varTrace"),
	}

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), metric.WithAttributes(otelAttrs...))
	}()

	// Increment request counter.
	c.requests.Add(ctx, 1, metric.WithAttributes(otelAttrs...))

	// Start a span for this request.
	ctx, span := c.cfg.Tracer.Start(ctx, "VarTracePost",
		trace.WithAttributes(otelAttrs...),
		clientSpanKind,
	)
	// Track stage for error reporting.
	var stage string
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			c.errors.Add(ctx, 1, metric.WithAttributes(otelAttrs...))
		}
		span.End()
	}()

	stage = "BuildURL"
	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/varTrace"
	uri.AddPathParts(u, pathParts[:]...)

	stage = "EncodeQueryParams"
	q := uri.NewQueryEncoder()
	{
		// Encode "policy" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Policy))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "commands" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "commands",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Commands))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "input" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "input",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Input.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "data" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "data",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Data.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "query" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "query",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(params.Query))
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	stage = "EncodeRequest"
	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			stage = "Security:BearerAuth"
			switch err := c.securityBearerAuth(ctx, "VarTracePost", r); {
			case err == nil: // if NO error
				satisfied[0] |= 1 << 0
			case errors.Is(err, ogenerrors.ErrSkipClientSecurity):
				// Skip this security.
			default:
				return res, errors.Wrap(err, "security \"BearerAuth\"")
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			return res, ogenerrors.ErrSecurityRequirementIsNotSatisfied
		}
	}

	stage = "SendRequest"
	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	stage = "DecodeResponse"
	result, err := decodeVarTracePostResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

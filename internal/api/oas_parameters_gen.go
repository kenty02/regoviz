// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// AstGetParams is parameters of GET /ast operation.
type AstGetParams struct {
	// The rego code to analyze.
	Policy string
}

func unpackAstGetParams(packed middleware.Parameters) (params AstGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	return params
}

func decodeAstGetParams(args [0]string, argsEscaped bool, r *http.Request) (params AstGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// AstPrettyGetParams is parameters of GET /astPretty operation.
type AstPrettyGetParams struct {
	// The rego code to analyze.
	Policy string
}

func unpackAstPrettyGetParams(packed middleware.Parameters) (params AstPrettyGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	return params
}

func decodeAstPrettyGetParams(args [0]string, argsEscaped bool, r *http.Request) (params AstPrettyGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// CallTreeAvailableEntrypointsGetParams is parameters of GET /callTree/availableEntrypoints operation.
type CallTreeAvailableEntrypointsGetParams struct {
	// The rego code to analyze.
	Policy string
}

func unpackCallTreeAvailableEntrypointsGetParams(packed middleware.Parameters) (params CallTreeAvailableEntrypointsGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	return params
}

func decodeCallTreeAvailableEntrypointsGetParams(args [0]string, argsEscaped bool, r *http.Request) (params CallTreeAvailableEntrypointsGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// CallTreeGetParams is parameters of GET /callTree operation.
type CallTreeGetParams struct {
	// The rego code to analyze.
	Policy string
	// The entrypoint rule to analyze.
	Entrypoint string
	// The input to policy.
	Input OptString
	// The data to policy.
	Data OptString
	// The query to policy.
	Query OptString
}

func unpackCallTreeGetParams(packed middleware.Parameters) (params CallTreeGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	{
		key := middleware.ParameterKey{
			Name: "entrypoint",
			In:   "query",
		}
		params.Entrypoint = packed[key].(string)
	}
	{
		key := middleware.ParameterKey{
			Name: "input",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Input = v.(OptString)
		}
	}
	{
		key := middleware.ParameterKey{
			Name: "data",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Data = v.(OptString)
		}
	}
	{
		key := middleware.ParameterKey{
			Name: "query",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Query = v.(OptString)
		}
	}
	return params
}

func decodeCallTreeGetParams(args [0]string, argsEscaped bool, r *http.Request) (params CallTreeGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: entrypoint.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "entrypoint",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Entrypoint = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "entrypoint",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: input.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "input",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotInputVal string
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					paramsDotInputVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Input.SetTo(paramsDotInputVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "input",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: data.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "data",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotDataVal string
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					paramsDotDataVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Data.SetTo(paramsDotDataVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "data",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: query.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "query",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotQueryVal string
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					paramsDotQueryVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Query.SetTo(paramsDotQueryVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "query",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// DepTreeTextGetParams is parameters of GET /depTreeText operation.
type DepTreeTextGetParams struct {
	// The rego code to analyze.
	Policy string
}

func unpackDepTreeTextGetParams(packed middleware.Parameters) (params DepTreeTextGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	return params
}

func decodeDepTreeTextGetParams(args [0]string, argsEscaped bool, r *http.Request) (params DepTreeTextGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// FlowchartGetParams is parameters of GET /flowchart operation.
type FlowchartGetParams struct {
	// The rego code to analyze.
	Policy string
	// Whether to return the editable flowchart mermaid url.
	Edit OptBool
}

func unpackFlowchartGetParams(packed middleware.Parameters) (params FlowchartGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	{
		key := middleware.ParameterKey{
			Name: "edit",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Edit = v.(OptBool)
		}
	}
	return params
}

func decodeFlowchartGetParams(args [0]string, argsEscaped bool, r *http.Request) (params FlowchartGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: edit.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "edit",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotEditVal bool
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToBool(val)
					if err != nil {
						return err
					}

					paramsDotEditVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Edit.SetTo(paramsDotEditVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "edit",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// IrGetParams is parameters of GET /ir operation.
type IrGetParams struct {
	// The rego code to analyze.
	Policy string
}

func unpackIrGetParams(packed middleware.Parameters) (params IrGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	return params
}

func decodeIrGetParams(args [0]string, argsEscaped bool, r *http.Request) (params IrGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// VarTracePostParams is parameters of POST /varTrace operation.
type VarTracePostParams struct {
	// The rego code to analyze.
	Policy string
	// The commands to analyze.
	Commands string
	// The input to policy.
	Input OptString
	// The data to policy.
	Data OptString
	// The query to policy.
	Query string
}

func unpackVarTracePostParams(packed middleware.Parameters) (params VarTracePostParams) {
	{
		key := middleware.ParameterKey{
			Name: "policy",
			In:   "query",
		}
		params.Policy = packed[key].(string)
	}
	{
		key := middleware.ParameterKey{
			Name: "commands",
			In:   "query",
		}
		params.Commands = packed[key].(string)
	}
	{
		key := middleware.ParameterKey{
			Name: "input",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Input = v.(OptString)
		}
	}
	{
		key := middleware.ParameterKey{
			Name: "data",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.Data = v.(OptString)
		}
	}
	{
		key := middleware.ParameterKey{
			Name: "query",
			In:   "query",
		}
		params.Query = packed[key].(string)
	}
	return params
}

func decodeVarTracePostParams(args [0]string, argsEscaped bool, r *http.Request) (params VarTracePostParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: policy.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "policy",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Policy = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "policy",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: commands.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "commands",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Commands = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "commands",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: input.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "input",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotInputVal string
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					paramsDotInputVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Input.SetTo(paramsDotInputVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "input",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: data.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "data",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotDataVal string
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToString(val)
					if err != nil {
						return err
					}

					paramsDotDataVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.Data.SetTo(paramsDotDataVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "data",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: query.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "query",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Query = c
				return nil
			}); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "query",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

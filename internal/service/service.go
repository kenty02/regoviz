package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/ir"
	"regoviz/internal/analyzer"
	"regoviz/internal/analyzer/astprint"
	"regoviz/internal/api"
	"regoviz/internal/samples"
	"regoviz/internal/utils"
	"strconv"
	"strings"
)

type regovizService struct{}

func (s *regovizService) CallTreeGet(ctx context.Context, params api.CallTreeGetParams) (*api.CallTreeGetOK, error) {
	// todo
	stub := api.CallTreeGetOK{
		Entrypoint: api.RuleParent{
			Name:    "i_am_entrypoint",
			UID:     utils.GenerateId(),
			Type:    api.RuleParentTypeParent,
			Default: "false",
			Children: []api.RuleParentChildrenItem{{
				Type: api.RuleChildRuleParentChildrenItem,
				RuleChild: api.RuleChild{
					Name:  "i_am_entrypoint_1",
					UID:   utils.GenerateId(),
					Type:  api.RuleChildTypeChild,
					Value: "",
					Statements: []api.RuleStatement{
						{
							Name: "foo == data.foo",
							UID:  utils.GenerateId(),
							Dependencies: []api.RuleStatementDependenciesItem{
								{
									Type: api.RuleParentRuleStatementDependenciesItem,
									RuleParent: api.RuleParent{
										Name:     "foo",
										UID:      utils.GenerateId(),
										Type:     api.RuleParentTypeParent,
										Default:  "false",
										Children: nil,
									},
								},
								{
									Type:   api.StringRuleStatementDependenciesItem,
									String: "data.foo",
								},
							},
						},
					},
				},
			},
			},
		},
	}
	return &stub, nil
}

func (s *regovizService) IrGet(ctx context.Context, params api.IrGetParams) (*api.IrGetOK, error) {
	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}
	policy, err := analyzer.PlanModuleAndGetIr(ctx, sample, false, true)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	if err := ir.Pretty(&buf, policy); err != nil {
		return nil, err
	}
	return &api.IrGetOK{Result: buf.String()}, nil
}

func (s *regovizService) FlowchartGet(ctx context.Context, params api.FlowchartGetParams) (*api.FlowchartGetOK, error) {
	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}
	plan, err := analyzer.PlanModuleAndGetIr(ctx, sample, false, true)
	if err != nil {
		return nil, err
	}

	mermaid := analyzer.GetMermaidFlowchart(plan)
	url, err := analyzer.GetMermaidUrl(mermaid, params.Edit.Or(false))
	if err != nil {
		return nil, err
	}

	// return modJson
	return &api.FlowchartGetOK{Result: url}, nil
}

func (s *regovizService) VarTracePost(_ context.Context, params api.VarTracePostParams) (*api.VarTracePostOK, error) {
	// convert params.Input  to map[string]interface{}
	var input map[string]interface{}
	if inputParam, ok := params.Input.Get(); ok {
		if err := json.Unmarshal([]byte(inputParam), &input); err != nil {
			return nil, err
		}
	}
	var data map[string]interface{}
	if dataParam, ok := params.Data.Get(); ok {
		if err := json.Unmarshal([]byte(dataParam), &data); err != nil {
			return nil, err
		}
	}
	query := params.Query
	var commands []interface{}
	commandsStr := params.Commands
	// convert params.Commands to []interface{}
	// "fixVar 8 role \"hogeeee\""
	// "showVars 8 role"
	if commandsStr != "" {
		commandsStrs := strings.Split(commandsStr, "\n")
		for _, commandStr := range commandsStrs {
			commandStr = strings.TrimSpace(commandStr)
			if commandStr == "" {
				continue
			}
			// skip # comments
			if strings.HasPrefix(commandStr, "#") {
				continue
			}
			commandStrs := strings.Split(commandStr, " ")
			if len(commandStrs) < 3 {
				return nil, fmt.Errorf("invalid command: %s", commandStr)
			}
			command := commandStrs[0]
			lineNum, err := strconv.Atoi(commandStrs[1])
			if err != nil {
				return nil, fmt.Errorf("invalid command: %s", commandStr)
			}
			varName := commandStrs[2]
			switch command {
			case "fixVar":
				varValue := strings.Join(commandStrs[3:], " ")
				commands = append(commands, analyzer.FixVarCommand{
					VarLineNum: lineNum,
					VarName:    varName,
					VarValue:   varValue,
				})
			case "showVars":
				commands = append(commands, analyzer.ShowVarsCommand{
					VarLineNum: lineNum,
					VarName:    varName,
				})
			default:
				return nil, fmt.Errorf("invalid command: %s", commandStr)
			}
		}
	}

	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}
	result, err := analyzer.DoVarTrace(sample, query, input, data, commands)
	if err != nil {
		return nil, err
	}
	return &api.VarTracePostOK{Result: result}, nil
}

func (s *regovizService) SamplesGet(_ context.Context) ([]api.Sample, error) {
	return samples.ListSamples("samples")
}

func (s *regovizService) AstGet(_ context.Context, params api.AstGetParams) (*api.AstGetOK, error) {
	//load sample
	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}
	mod, err := analyzer.CompileModuleStringToAst(sample)

	if err != nil {
		return nil, err
	}

	// serialize module to modJson
	modJson, err := json.Marshal(mod)

	if err != nil {
		return nil, err
	}

	// return modJson
	return &api.AstGetOK{Result: string([]byte(modJson))}, nil
}

func (s *regovizService) AstPrettyGet(_ context.Context, params api.AstPrettyGetParams) (*api.AstPrettyGetOK, error) {
	//load sample
	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}
	// compile module
	mod, err := analyzer.CompileModuleStringToAst(sample)

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}

	err = astprint.Pretty(&buf, mod)

	if err != nil {
		return nil, err
	}

	// return modJson
	return &api.AstPrettyGetOK{Result: buf.String()}, nil
}

func (s *regovizService) DepTreeTextGet(ctx context.Context, params api.DepTreeTextGetParams) (*api.DepTreeTextGetOK, error) {
	//// compile module
	//mod, err := compileRego(params.Module)
	//
	//if err != nil {
	//	return nil, err
	//}
	//// string builder
	//var sb strings.Builder
	//
	//// iterate each rule
	//for _, rule := range mod.Rules {
	//	// iterate each expression
	//	for _, expr := range rule.Body {
	//
	//		// iterate each term
	//		terms, ok := expr.Terms.([]*ast.Term)
	//		if !ok {
	//			terms = []*ast.Term{expr.Terms.(*ast.Term)}
	//		}
	//		for _, term := range terms {
	//			// check term value is ast.Ref array
	//			var termValue ast.Ref
	//			ok := false
	//			if termValue, ok = term.Value.(ast.Ref); !ok {
	//				continue
	//			}
	//			// iterate each value
	//			for _, value := range termValue {
	//				// iterate each referencedRule
	//				for _, referencedRule := range mod.Rules {
	//					// check referencedRule head name is equal to value
	//					// print both hand side of this if
	//					ruleName, err := ast.InterfaceToValue(referencedRule.Head.Name.String())
	//					if err != nil {
	//						return nil, err
	//					}
	//					if value.Value.Compare(ruleName) == 0 {
	//						// write to string builder
	//						sb.WriteString(fmt.Sprintf("%s -> %s\n", rule.Head.Name, referencedRule.Head.Name))
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//resultString := sb.String()
	//return &api.DepTreeTextGetOK{Result: resultString}, nil

	sample, err := samples.ReadSample(params.SampleName, "samples")
	if err != nil {
		return nil, err
	}

	plan, err := analyzer.PlanModuleAndGetIr(ctx, sample, false, true)
	if err != nil {
		return nil, err
	}
	treeMap := analyzer.GetDepTreePretty(plan)
	return &api.DepTreeTextGetOK{Result: treeMap}, nil
}

func NewService() api.Handler {
	return &regovizService{}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	vizrego "vizrego-poc/vizrego"
)

type vizRegoService struct {
	id  int64
	mux sync.Mutex
}

func (p *vizRegoService) FlowchartGet(ctx context.Context, params vizrego.FlowchartGetParams) (*vizrego.FlowchartGetOK, error) {
	sample, err := readSample(params.SampleName)
	if err != nil {
		return nil, err
	}
	plan, err := plan(ctx, sample, false)
	if err != nil {
		return nil, err
	}

	mermaid := getMermaidFlowchart(plan)
	url := getMermaidUrl(mermaid)

	// return modJson
	return &vizrego.FlowchartGetOK{Result: url}, nil
}

func (p *vizRegoService) VarTracePost(ctx context.Context, params vizrego.VarTracePostParams) (*vizrego.VarTracePostOK, error) {
	// convert params.Input  to map[string]interface{}
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(params.Input), &input); err != nil {
		return nil, err
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
				varValue := commandStrs[3]
				commands = append(commands, FixVarCommand{
					varLineNum: lineNum,
					varName:    varName,
					varValue:   varValue,
				})
			case "showVars":
				commands = append(commands, ShowVarsCommand{
					varLineNum: lineNum,
					varName:    varName,
				})
			default:
				return nil, fmt.Errorf("invalid command: %s", commandStr)
			}
		}
	}

	sample, err := readSample(params.SampleName)
	if err != nil {
		return nil, err
	}
	result, err := regoVarTrace(sample, query, input, commands)
	if err != nil {
		return nil, err
	}
	return &vizrego.VarTracePostOK{Result: result}, nil
}

func NewService() vizrego.Handler {
	return &vizRegoService{}
}

func (p *vizRegoService) SamplesGet(ctx context.Context) ([]vizrego.Sample, error) {
	// find .rego files in samples directory and return them
	files, err := os.ReadDir("samples")
	if err != nil {
		return nil, err
	}
	var samples []vizrego.Sample
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".rego") {
			// read content
			content, err := os.ReadFile(fmt.Sprintf("samples/%s", file.Name()))
			if err != nil {
				return nil, err
			}
			samples = append(samples, vizrego.Sample{
				FileName: file.Name(),
				Content:  string(content),
			})
		}
	}
	return samples, nil
}

func (p *vizRegoService) RulesGet(ctx context.Context) ([]vizrego.Rule, error) {
	// not implemented
	return nil, error(nil)
}

func (p *vizRegoService) AstGet(ctx context.Context, params vizrego.AstGetParams) (*vizrego.AstGetOK, error) {
	// compile module
	mod, err := compileRego(params.Module)

	if err != nil {
		return nil, err
	}

	// serialize module to modJson
	modJson, err := json.Marshal(mod)

	if err != nil {
		return nil, err
	}

	// return modJson
	return &vizrego.AstGetOK{Result: string([]byte(modJson))}, nil
}

func (p *vizRegoService) DepTreeTextGet(ctx context.Context, params vizrego.DepTreeTextGetParams) (*vizrego.DepTreeTextGetOK, error) {
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
	//return &vizrego.DepTreeTextGetOK{Result: resultString}, nil

	plan, err := plan(ctx, params.Module, false)
	if err != nil {
		return nil, err
	}
	treeMap := getDepTreePretty(plan)
	return &vizrego.DepTreeTextGetOK{Result: treeMap}, nil
}

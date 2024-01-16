package evalTrace

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/topdown/print"
	"os"
	"regoviz/internal/analyzer"
	"regoviz/internal/analyzer/regomod"
	"regoviz/internal/api"
)

//var beginRuleParentRegex *regexp.Regexp = regexp.MustCompile(`^begin_rule_parent (.+)$`)
//var endRuleParentRegex *regexp.Regexp = regexp.MustCompile(`^end_rule_parent (.+) (.*)$`)
//var beginRuleChildRegex *regexp.Regexp = regexp.MustCompile(`^begin_rule_child (.+) (\d+)$`)
//var endRuleChildRegex *regexp.Regexp = regexp.MustCompile(`^end_rule_child (.+) (\d+)$`)
//var beginRuleStatementRegex *regexp.Regexp = regexp.MustCompile(`^begin_rule_statement (\d+)$`)
//var endRuleStatementRegex *regexp.Regexp = regexp.MustCompile(`^end_rule_statement (\d+)$`)
//var ruleStatementVarRegex *regexp.Regexp = regexp.MustCompile(`^rule_statement_var (.+) (.*)$`)

// depths is a helper for computing the depth of an event. Events within the
// same query all have the same depth. The depth of query is
// depth(parent(query))+1.
type depths map[uint64]int

func (ds depths) GetOrSet(qid uint64, pqid uint64) int {
	depth := ds[qid]
	if depth == 0 {
		depth = ds[pqid]
		depth++
		ds[qid] = depth
	}
	return depth
}

type printEvent struct {
	location *ast.Location
	message  string
}

type step struct {
	// index starts at 1
	index int
	// either printEvent or topdown.Event and its depth is provided
	printEvent *printEvent

	depth int
	event *topdown.Event
}

type tracer struct {
	steps         *[]*step
	nodeLocToStep map[string]*step
	depths        *depths
	filter        func(*topdown.Event) bool
}

func (t tracer) Enabled() bool {
	return true
}

func (t tracer) TraceEvent(event topdown.Event) {
	if t.filter != nil && !t.filter(&event) {
		return
	}
	step := &step{
		index: len(*t.steps) + 1,
		depth: t.depths.GetOrSet(event.QueryID, event.ParentID),
		event: &event,
	}
	*t.steps = append(*t.steps, step)
	if event.Node != nil {
		if loc := event.Node.Loc(); loc != nil {
			t.nodeLocToStep[fmt.Sprintf("%v:%v", loc.Row, loc.Col)] = step
		}
	}
}

func (t tracer) Config() topdown.TraceConfig {
	return topdown.TraceConfig{PlugLocalVars: true}
}

func newTracer() tracer {
	return tracer{
		steps:         &[]*step{},
		nodeLocToStep: map[string]*step{},
		depths:        &depths{},
	}
}

func (t tracer) Print(context print.Context, s string) error {
	step := &step{
		index:      len(*t.steps) + 1,
		printEvent: &printEvent{location: context.Location, message: s},
	}
	*t.steps = append(*t.steps, step)
	return nil
}

func locationRowColEqual(loc api.NodeLocation, other *ast.Location) bool {
	return loc.Row == other.Row && loc.Col == other.Col
}

func DoEvalTrace(rego, query string, input, data map[string]interface{}, callTree *api.RuleParent, callTreeNodes []interface{}, commands []interface{}) ([]api.EvalStep, error) {
	modifiedRego, _, err := regomod.Apply(regomod.Opts{
		Rego:                  rego,
		RuleParentTrace:       false,
		RuleChildTrace:        false,
		RuleStatementTrace:    false,
		RuleStatementVarTrace: true,
		// todo support below
		RuleStatementVarAllTraceTarget: nil,
		RuleStatementVarFixes:          nil,
	})
	if err != nil {
		return nil, err
	}

	tracer := newTracer()
	tracer.filter = func(event *topdown.Event) bool {
		return event.Op != topdown.UnifyOp
	}

	_, err = analyzer.EvalRegoWithPrintAndTrace(modifiedRego, query, input, data, tracer, tracer)
	if err != nil {
		return nil, err
	}

	//debug
	//var events []*topdown.Event
	//for _, step := range *tracer.steps {
	//	if step.printEvent == nil {
	//		events = append(events, step.event)
	//	}
	//}
	//var buf bytes.Buffer
	//buf.Grow(1024 * 1024)
	//topdown.PrettyTraceWithLocation(&buf, events)
	//fmt.Println(buf.String())

	isEventInQuery := func(event *topdown.Event) bool {
		return event.Location.File == ""
	}

	var evalSteps []api.EvalStep
	var previousTargetNodeUid string
	for _, step := range *tracer.steps {
		if step.printEvent != nil {
			evalStep := &api.EvalStep{
				Index:         step.index,
				Message:       "(print) " + step.printEvent.message,
				TargetNodeUid: previousTargetNodeUid,
			}
			evalSteps = append(evalSteps, *evalStep)
			continue
		}
		matched := false
		switch step.event.Op {
		case topdown.EnterOp, topdown.ExitOp, topdown.FailOp, topdown.RedoOp, topdown.EvalOp:
			for _, callTreeNode := range callTreeNodes {
				var callTreeNodeLocation *api.OptNodeLocation
				var callTreeNodeUid string
				var callTreeNodeLocationText string
				switch callTreeNode := callTreeNode.(type) {
				case *api.RuleChild:
					callTreeNodeLocation = &callTreeNode.Location
					callTreeNodeUid = callTreeNode.UID
					rule, ok := step.event.Node.(*ast.Rule)
					if !ok {
						continue
					}
					callTreeNodeLocationText = fmt.Sprintf("RULE `%s` at row %d", rule.Head.Name, rule.Head.Loc().Row)
				case *api.RuleStatement:
					callTreeNodeLocation = &callTreeNode.Location
					callTreeNodeUid = callTreeNode.UID
					expr, ok := step.event.Node.(*ast.Expr)
					if !ok {
						continue
					}
					callTreeNodeLocationText = fmt.Sprintf("STATEMENT `%s` at row %d", expr.Location.Text, expr.Location.Row)
				default:
					continue
				}
				if !callTreeNodeLocation.Set {
					return nil, fmt.Errorf("location not set for %T", callTreeNode)
				}
				if isEventInQuery(step.event) || !locationRowColEqual(callTreeNodeLocation.Value, step.event.Node.Loc()) {
					continue
				}
				message := ""
				switch step.event.Op {
				case topdown.EnterOp:
					message = fmt.Sprintf("Evaluating %s", callTreeNodeLocationText)
				case topdown.ExitOp:
					message = fmt.Sprintf("Evaluated to be TRUTHY: %s", callTreeNodeLocationText)
				case topdown.RedoOp:
					message = fmt.Sprintf("Re-evaluating %s", callTreeNodeLocationText)
				case topdown.EvalOp:
					message = fmt.Sprintf("Evaluating %s", callTreeNodeLocationText)
				case topdown.FailOp:
					message = fmt.Sprintf("Evaluated to be FALSY: %s", callTreeNodeLocationText)
				default:
					panic("unexpected")
				}
				evalStep := &api.EvalStep{
					Index:         step.index,
					Message:       message,
					TargetNodeUid: callTreeNodeUid,
				}
				previousTargetNodeUid = callTreeNodeUid
				evalSteps = append(evalSteps, *evalStep)
				matched = true
			}
		case topdown.NoteOp:
			matched = true
			evalStep := &api.EvalStep{
				Index:         step.index,
				Message:       fmt.Sprintf("(trace) %s", step.event.Message),
				TargetNodeUid: previousTargetNodeUid,
			}
			evalSteps = append(evalSteps, *evalStep)
		case topdown.IndexOp:
			for _, callTreeNode := range callTreeNodes {
				ruleParent, ok := callTreeNode.(*api.RuleParent)
				if !ok {
					continue
				}
				if ruleParent.Ref != step.event.Ref.String() {
					continue
				}
				evalStep := &api.EvalStep{
					Index:         step.index,
					Message:       "Looking up " + step.event.Ref.String(),
					TargetNodeUid: ruleParent.UID,
				}
				previousTargetNodeUid = ruleParent.UID
				evalSteps = append(evalSteps, *evalStep)
				matched = true
			}
		default:
			//unimplemented
		}
		addUnsupported := os.Getenv("ADD_UNSUPPORTED_STEP") == "true"
		if !matched && addUnsupported {
			evalStep := &api.EvalStep{
				Index:         step.index,
				Message:       fmt.Sprintf("(unsupported at %s) %s", fmt.Sprintf("%s %d:%d `%s`", step.event.Location.File, step.event.Location.Row, step.event.Location.Col, string(step.event.Location.Text)), step.event.String()),
				TargetNodeUid: "",
			}
			evalSteps = append(evalSteps, *evalStep)
		}
	}
	return evalSteps, nil
}

package evalTrace

import (
	"bytes"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/topdown/print"
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
}

func (t tracer) Enabled() bool {
	return true
}

func (t tracer) TraceEvent(event topdown.Event) {
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
	_, err = analyzer.EvalRegoWithPrintAndTrace(modifiedRego, query, input, data, tracer, tracer)
	if err != nil {
		return nil, err
	}

	//debug
	var events []*topdown.Event
	for _, step := range *tracer.steps {
		if step.printEvent == nil {
			events = append(events, step.event)
		}
	}
	var buf bytes.Buffer
	topdown.PrettyTraceWithLocation(&buf, events)
	//fmt.Println(buf.String())

	isEventInQuery := func(event *topdown.Event) bool {
		return event.Location.File == ""
	}

	var evalSteps []api.EvalStep
	var previousNonPrintStep *step
	for _, step := range *tracer.steps {
		if step.printEvent != nil {
			if previousNonPrintStep == nil {
				return nil, fmt.Errorf("unexpected print event")
			}
			evalStep := &api.EvalStep{
				Index:         step.index,
				Message:       "(print) " + step.printEvent.message,
				TargetNodeUid: "TODO",
			}
			evalSteps = append(evalSteps, *evalStep)
			continue
		}
		switch step.event.Op {
		case topdown.EnterOp:
			match := 0
			for _, callTreeNode := range callTreeNodes {
				ruleChild, ok := callTreeNode.(*api.RuleChild)
				if !ok {
					continue
				}
				if !ruleChild.Location.Set {
					return nil, fmt.Errorf("location not set for %T", ruleChild)
				}
				if !isEventInQuery(step.event) && locationRowColEqual(ruleChild.Location.Value, step.event.Node.Loc()) {
					evalStep := &api.EvalStep{
						Index:         step.index,
						Message:       step.event.String(),
						TargetNodeUid: ruleChild.UID,
					}
					evalSteps = append(evalSteps, *evalStep)
					match++
				}
			}
			// todo check match
		case topdown.EvalOp:
			if isEventInQuery(step.event) {
				continue
			}
			match := 0
			for _, callTreeNode := range callTreeNodes {
				ruleStatement, ok := callTreeNode.(*api.RuleStatement)
				if !ok {
					continue
				}
				expr, ok := step.event.Node.(*ast.Expr)
				if !ok {
					panic("unexpected")
				}
				if !ruleStatement.Location.Set {
					return nil, fmt.Errorf("location not set for %T", ruleStatement)
				}
				if locationRowColEqual(ruleStatement.Location.Value, expr.Loc()) {
					evalStep := &api.EvalStep{
						Index:         step.index,
						Message:       step.event.String(),
						TargetNodeUid: ruleStatement.UID,
					}
					evalSteps = append(evalSteps, *evalStep)
					match++
				}
			}
			// todo check match
		case topdown.NoteOp:
			// todo
		case topdown.IndexOp:
			for _, callTreeNode := range callTreeNodes {
				match := 0
				ruleParent, ok := callTreeNode.(*api.RuleParent)
				if !ok {
					continue
				}
				if ruleParent.Ref == step.event.Ref.String() {
					evalStep := &api.EvalStep{
						Index:         step.index,
						Message:       step.event.String(),
						TargetNodeUid: ruleParent.UID,
					}
					evalSteps = append(evalSteps, *evalStep)
					match++
				}
				// todo check match
			}
		default:
			//unimplemented
		}
		previousNonPrintStep = step
	}
	return evalSteps, nil
}

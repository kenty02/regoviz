package regomod

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"reflect"
	"regexp"
	"regoviz/internal/analyzer"
	"strings"
)

type (
	ValTarget struct {
		Row int
		Val string
	}
	ValFixTarget struct {
		ValTarget
		NewValue string
	}
	// Opts is the options for Apply. Only RuleStatementTrace is supported for now.
	Opts struct {
		Rego                           string
		RuleParentTrace                bool
		RuleChildTrace                 bool
		RuleStatementTrace             bool
		RuleStatementVarTrace          bool
		RuleStatementVarAllTraceTarget *ValTarget
		RuleStatementVarFixes          []*ValFixTarget
	}
	VarTraceFoundVars struct {
		Rule *ast.Rule
		Expr *ast.Expr
		Vars ast.VarSet
	}
)

func replaceTerm(termMap map[*ast.Term]string, target string) (string, error) {
	rows := strings.Split(strings.ReplaceAll(target, "\r", ""), "\n")
	modifiedRowNums := map[int]bool{}
	for term, replacement := range termMap {
		if rowModified, ok := modifiedRowNums[term.Location.Row]; ok && rowModified {
			return "", fmt.Errorf("multiple term modifications on row %d (not supported yet)", term.Location.Row)
		}
		row := term.Location.Row
		col := term.Location.Col

		rowString := rows[row-1]

		newRowString := rowString[:col-1] + replacement + rowString[col-1+len(term.Location.Text):]
		rows[row-1] = newRowString

		modifiedRowNums[term.Location.Row] = true
	}
	return strings.Join(rows, "\n"), nil
}

func replaceExpr(exprMap map[*ast.Expr]string, target string, rule *ast.Rule) (string, error) {
	target = addBracketToRuleBody(rule, target)
	rows := strings.Split(strings.ReplaceAll(target, "\r", ""), "\n")
	oldnewForRows := map[int][]string{}
	for expr, replacement := range exprMap {
		oldnewForRows[expr.Location.Row] = append(oldnewForRows[expr.Location.Row], string(expr.Location.Text))
		oldnewForRows[expr.Location.Row] = append(oldnewForRows[expr.Location.Row], replacement)
	}
	for row, oldnew := range oldnewForRows {
		rowString := rows[row-1]
		replacer := strings.NewReplacer(oldnew...)
		newRowString := replacer.Replace(rowString)
		rows[row-1] = newRowString
	}
	return strings.Join(rows, "\n"), nil
}

var ruleRegex = regexp.MustCompile(`^(.+) if (.+)$`)

// rule must use "if" keyword in order to work
// rules must not have exactly same body
func addBracketToRuleBody(rule *ast.Rule, target string) string {
	text := rule.Location.Text
	match := ruleRegex.FindSubmatch(text)
	if match == nil {
		return target
	}
	if strings.HasPrefix(string(match[1]), "{") {
		return target
	}

	// It's very hackey, but by not putting a space between the if and the {, we make sure that the start of the first expr doesn't change.
	newText := fmt.Sprintf("%s if{%s}", match[1], match[2])
	return strings.Replace(target, string(text), newText, 1)
}
func Apply(o Opts) (string, []VarTraceFoundVars, error) {
	module, moduleCompiler, err := analyzer.CompileModuleStringToAst(o.Rego, false, true)
	if err != nil {
		return "", nil, err
	}
	modifiedRego := o.Rego
	maybeGenerateId := func() string {
		return "0"
		//return utils.GenerateId()
	}

	rules := module.Rules
	if len(rules) == 0 {
		return "", nil, fmt.Errorf("no rules found")
	}
	var varTraceFoundVarsList []VarTraceFoundVars

	ruleStatementVarAllTraceTargetRemains := o.RuleStatementVarAllTraceTarget != nil
	for _, rule := range rules {
		if rule.Default {
			continue
		}
		// rule.generatedBody == true
		if reflect.ValueOf(rule).Elem().FieldByName("generatedBody").Bool() {
			continue
		}
		modifiedExprs := make([]string, len(rule.Body))
		exprToModifiedExpr := make(map[*ast.Expr]string)
		for i, expr := range rule.Body {
			modifiedExprs[i] = string(expr.Location.Text)
			// first statement
			if o.RuleChildTrace && i == 0 {
				modifiedExprs[i] = fmt.Sprintf(`trace("begin_rule_child %s %d");%s`, rule.Head.Name.String(), rule.Head.Location.Row, modifiedExprs[i])
			}
			// statement begin
			if o.RuleStatementTrace {
				modifiedExprs[i] =
					fmt.Sprintf(`trace("begin_rule_statement %d");%s`,
						i, modifiedExprs[i])
			}

			// statement var trace & fix
			if o.RuleStatementVarTrace {
				vars := ast.VarSet{}
				processTerm := func(term *ast.Term) (bool, error) {
					v, ok := term.Value.(ast.Var)
					if !ok {
						return false, nil
					}
					if v.IsGenerated() {
						after, ok := moduleCompiler.RewrittenVars[v]
						if ok && !after.IsGenerated() {
							v = after
						} else {
							return false, nil
						}
					}
					// fix
					for _, fix := range o.RuleStatementVarFixes {
						if fix.Row == expr.Location.Row && fix.Val == v.String() {
							newValName := v.String() + "__fix_" + maybeGenerateId()
							modified, err := replaceTerm(map[*ast.Term]string{term: newValName}, modifiedExprs[i])
							if err != nil {
								return false, err
							}
							modifiedExprs[i] = fmt.Sprintf(`%s;%s:=%s`, modified, fix.Val, fix.NewValue)
						}
					}

					// trace
					if o.RuleStatementVarTrace {
						maybeFalseStatement := ""
						if ruleStatementVarAllTraceTargetRemains {
							target := o.RuleStatementVarAllTraceTarget
							if target.Row == expr.Location.Row && target.Val == v.String() {
								maybeFalseStatement = ";false"
								ruleStatementVarAllTraceTargetRemains = false
							}
						}
						modifiedExprs[i] =
							fmt.Sprintf(`%s;trace(sprintf("rule_statement_var %s %%v", [%s]))%s`,
								modifiedExprs[i], v, v, maybeFalseStatement)
						return true, nil
					}
					return false, nil
				}
				switch t := expr.Terms.(type) {
				case []*ast.Term:
					for _, term := range t {
						needBreak, err := processTerm(term)
						if err != nil {
							return "", nil, err
						}
						if needBreak {
							break
						}
					}
				case *ast.Term:
					_, err := processTerm(t)
					if err != nil {
						return "", nil, err
					}
				default:
					return "", nil, fmt.Errorf("unsupported term type: %T", t)
				}
				if len(vars) > 0 {
					varTraceFoundVarsList = append(varTraceFoundVarsList, VarTraceFoundVars{
						Rule: rule,
						Expr: expr,
						Vars: vars,
					})
				}
			}

			// statement end
			if o.RuleStatementTrace {
				modifiedExprs[i] =
					fmt.Sprintf(`%s;trace("end_rule_statement %d")`,
						modifiedExprs[i], i)
			}

			// last statement
			if o.RuleStatementTrace && i == len(rule.Body)-1 {
				modifiedExprs[i] = fmt.Sprintf(`%s;trace("end_rule_child %s %d")`, modifiedExprs[i], rule.Head.Name.String(), rule.Head.Location.Row)
			}

			exprToModifiedExpr[expr] = modifiedExprs[i]
		}
		modifiedRego, err = replaceExpr(exprToModifiedExpr, modifiedRego, rule)
		if err != nil {
			return "", nil, err
		}
	}
	if ruleStatementVarAllTraceTargetRemains {
		return "", nil, fmt.Errorf("ruleStatementVarAllTraceTarget remained unprocessed")
	}

	newRulesAndBody := make(map[*ast.Rule]string)
	termsToReplace := make(map[*ast.Term]string)
	if o.RuleParentTrace {
		firstRule := rules[0]
		originRuleName := firstRule.Head.Name.String()

		// create wrapping rule
		originNewRuleName := firstRule.Head.Name.String() + "__origin_" + maybeGenerateId()
		// todo support args
		newRuleBody := fmt.Sprintf(`trace("begin_rule_parent %s");%s;trace("end_rule_parent %s")`, originRuleName, originNewRuleName, originRuleName)
		newRule := &ast.Rule{
			Default:  false,
			Head:     ast.NewHead(ast.Var(originRuleName), firstRule.Head.Args...),
			Body:     ast.MustParseBody("\"REPLACE_ME\""),
			Else:     nil,
			Location: nil,
			Module:   nil,
		}
		//firstRule.Module.Rules = append(firstRule.Module.Rules, newRule)
		newRulesAndBody[newRule] = newRuleBody

		// change original rule name to new rule name
		for _, rule := range rules {
			if rule.Head.Name.Equal(firstRule.Head.Name) {
				termsToReplace[rule.Head.Reference[0]] = originNewRuleName
			}
		}
	}

	modifiedRego, err = replaceTerm(termsToReplace, modifiedRego)
	if err != nil {
		return "", nil, err
	}
	for rule, body := range newRulesAndBody {
		ruleString := rule.String()
		ruleString = strings.ReplaceAll(ruleString, "\"REPLACE_ME\"", body)
		modifiedRego = modifiedRego + "\n" + ruleString
	}

	return modifiedRego, varTraceFoundVarsList, err
}

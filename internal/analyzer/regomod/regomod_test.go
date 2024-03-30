package regomod

import (
	"github.com/google/go-cmp/cmp"
	"github.com/open-policy-agent/opa/ast"
	"regoviz/internal/analyzer"
	"testing"
)

func term(s string, row, col int) *ast.Term {
	return ast.MustParseTerm("allow").
		SetLocation(&ast.Location{Text: []byte(s), Row: row, Col: col, Offset: len([]byte(s))})
}
func TestReplaceTerm(t *testing.T) {
	tests := []struct {
		description string
		input       map[*ast.Term]string
		target      string
		expected    string
	}{
		{
			description: "test1",
			input: map[*ast.Term]string{
				term("allow", 2, 1): "allow__origin_0",
			},
			target: `package test
allow := true`,
			expected: `package test
allow__origin_0 := true`,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual, err := replaceTerm(test.input, test.target)
			if err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(test.expected, actual); d != "" {
				t.Errorf("Result and expected objects differ for \"%s\": (-want +got)\n%s", test.description, d)
			}
		})
	}
}

type testData struct {
	description string
	input       Opts
	expected    string
}

func TestApply(t *testing.T) {
	tests := []testData{
		{
			description: "test1",
			input: Opts{
				RuleParentTrace: true,
				Rego: `package test
allow := true`,
			},
			expected: `package test
allow__origin_0 := true
allow { trace("begin_rule_parent allow");allow__origin_0;trace("end_rule_parent allow") }`,
		},
		// todo: pass below tests
		//		{
		//			description: "test2",
		//			input: Opts{
		//				RuleChildTrace: true,
		//				Rego: `package test
		//allow := true`,
		//			},
		//			expected: `package test
		//allow := true { print("begin_rule_child allow 2");print("end_rule_child allow 2") }`,
		//		},
		//		{
		//			description: "test3",
		//			input: Opts{
		//				Rego: `package test
		//import future.keywords.if
		//import future.keywords.in
		//allow if {
		//	some role in input.roles
		//	input.foo[_] == role
		//	input.bar == "baz"
		//}`,
		//				RuleParentTrace:                true,
		//				RuleChildTrace:                 true,
		//				RuleStatementTrace:             true,
		//				RuleStatementVarTrace:          true,
		//				RuleStatementVarAllTraceTarget: nil,
		//				RuleStatementVarFixes:          nil,
		//			},
		//			expected: `package test
		//import future.keywords.if
		//import future.keywords.in
		//allow__origin_0 if {
		//	print("begin_rule_statement 0");print("begin_rule_child allow 4");__local2__ = input.roles[__local1__];print("rule_statement_var __local2__ role", __local2__);print("end_rule_statement 0")
		//	print("begin_rule_statement 1");input.foo[_] = __local2__;print("rule_statement_var __local2__ role", __local2__);print("end_rule_statement 1")
		//	print("begin_rule_statement 2");input.bar = "baz";print("end_rule_statement 2");print("end_rule_child allow 4")
		//}
		//allow { print("begin_rule_parent allow");print("end_rule_parent allow", allow__origin_0) }`,
		//		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// todo check []VarTraceFoundVars
			actual, _, err := Apply(test.input)
			if err != nil {
				t.Fatal(err)
			}
			// check if the modified rego is compilable
			_, _, err = analyzer.CompileModuleStringToAst(actual, true, true)
			if err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(test.expected, actual); d != "" {
				t.Errorf("Result and expected objects differ for \"%s\": (-want +got)\n%s", test.description, d)
			}
		})
	}
}

package analyzer

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"regoviz/internal/api"
	"testing"
)

//go:embed testdata/rbac.rego
var rbacRego string

func TestCompileStringToAst(t *testing.T) {
	rego := rbacRego
	mod, _, err := CompileModuleStringToAst(rego)
	if err != nil {
		t.Fatal(err)
	}
	if mod.Rules[0].Head.Name != "allow" {
		t.Fatal("allow rule is not found")
	}
	if mod.Rules[2].Body.String() != "__local2__ = data.app.rbac.user_is_granted[__local1__]; input.action = __local2__.action; input.type = __local2__.type" {
		t.Fatal("user_is_granted body is wrong")
	}
}

func TestPlanModuleAndGetIrWithMetadata(t *testing.T) {
	ctx := context.Background()
	rego := `package test

import data.a

default allow = false

# METADATA
# entrypoint: true
allow {
	a[_] = input
}
`
	_, err := PlanModuleAndGetIr(ctx, rego, false, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlanWithoutMetadataWithoutMetadata(t *testing.T) {
	ctx := context.Background()
	rego := `package test

import data.a

default allow = false

allow {
	a[_] = input
}
`
	_, err := PlanModuleAndGetIr(ctx, rego, false, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlanModuleAndGetIrWithoutMetadataPackageContainsDot(t *testing.T) {
	ctx := context.Background()
	rego := `package te.st

import data.a

default allow = false

allow {
	a[_] = input
}
`
	_, err := PlanModuleAndGetIr(ctx, rego, false, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlanAsText(t *testing.T) {
	ctx := context.Background()
	rego := `package test

import data.a

default allow = false

allow {
	a[_] = input
}
`

	plan, err := PlanAsText(ctx, rego, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(plan)
	if len(plan) < 10 {
		t.Fatal("plan is too short")
	}

}

func TestGetDepTreeMap(t *testing.T) {
	ctx := context.Background()

	plan, err := PlanModuleAndGetIr(ctx, rbacRego, false, true)
	if err != nil {
		t.Fatal(err)
	}

	treeMap := GetDepTreeMap(plan)

	fmt.Println(treeMap)

	if len(treeMap) < 1 {
		t.Fatal("treeMap is too short")
	}
}

func TestGetMermaidFlowchart(t *testing.T) {
	ctx := context.Background()

	plan, err := PlanModuleAndGetIr(ctx, rbacRego, false, true)
	if err != nil {
		t.Fatal(err)
	}

	mermaid := GetMermaidFlowchart(plan)

	url, err := GetMermaidUrl(mermaid, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
	fmt.Println(mermaid)

	if len(mermaid) < 1 {
		t.Fatal("mermaid is too short")
	}
}

func TestEvalRegoWithPrint(t *testing.T) {
	rego := `package example

			rule_containing_print_call {
				print("input.foo is:", input.foo, "and input.bar is:", input.bar)
			}
		`
	query := "data.example.rule_containing_print_call"
	input := map[string]interface{}{
		"foo": 7}
	rs, buf, err := evalRegoWithPrint(rego, query, input, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rs.Allowed())
	fmt.Println(buf)
	if len(buf) < 1 {
		t.Fatal("buf is too short")
	}
}

func TestEvalRegoWithHTTPSendFails(t *testing.T) {
	rego := `package example

			allow {
				http.send({"method": "get", "url": "https://www.google.com"})
			}
		`
	query := "data.example.allow"
	_, _, err := evalRegoWithPrint(rego, query, nil, nil)
	if err == nil {
		t.Fatal("error should occur")
	}
}

// テストケースを実行するヘルパー関数
func testInjectCode(t *testing.T, originalCode string, injections []CodeInject, expectedCode string) {
	t.Helper()
	result := injectCode(originalCode, injections)
	if result != expectedCode {
		t.Errorf("Expected result does not match actual result.\nExpected:\n%s\nActual:\n%s", expectedCode, result)
	}
}

// injectCode関数のテスト
func TestInjectCode(t *testing.T) {
	originalCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	injections := []CodeInject{
		{1, " // This is the main package", false},
		{3, " // Importing the fmt package", false},
		{5, " // Entry point of the program", false},
	}

	expectedCode := `package main // This is the main package

import "fmt" // Importing the fmt package

func main() { // Entry point of the program
	fmt.Println("Hello, World!")
}`

	testInjectCode(t, originalCode, injections, expectedCode)
}

// 行番号が範囲外の場合のテスト
func TestInjectCodeOutOfRange(t *testing.T) {
	originalCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	// 行番号がコードの行数を超えている
	injections := []CodeInject{
		{10, " // This line number is out of range", false},
	}

	// 期待されるコードは元のコードと同じであるべきです。
	expectedCode := originalCode

	testInjectCode(t, originalCode, injections, expectedCode)
}

// 挿入すべきコードが空の場合のテスト
func TestInjectCodeEmptyInjection(t *testing.T) {
	originalCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	// 空のコードを挿入する
	injections := []CodeInject{
		{3, "", false},
	}

	// 期待されるコードは元のコードと同じであるべきです。
	expectedCode := originalCode

	testInjectCode(t, originalCode, injections, expectedCode)
}

func TestEvalRegoWithPrintInjected(t *testing.T) {
	rego := `package example
import future.keywords.in
import future.keywords.if

roles := ["admin", "user"]

allow {
	some role in roles
	# check if the user has the valid role
	input.role == role
}
		`
	rego = injectCode(rego, []CodeInject{
		{8, ";print(\"role =\", role);false", false},
	})
	fmt.Println(rego)
	query := "data.example.allow"
	input := map[string]interface{}{"role": "admin"}
	rs, buf, err := evalRegoWithPrint(rego, query, input, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rs.Allowed())
	fmt.Println(buf)
	if buf != "role = admin\nrole = user\n" {
		t.Fatal("buf is wrong")
	}
}

func TestVarTrace(t *testing.T) {
	rego := `package example
import future.keywords.in
import future.keywords.if

roles := ["admin", "user"]

allow {
	some role in roles
	# check if the user has the valid role
	input.role == role
}
`
	query := "data.example.allow"
	input := map[string]interface{}{"role": "admin"}
	commands := []interface{}{
		FixVarCommand{
			VarLineNum: 8,
			VarName:    "role",
			VarValue:   "\"hoge\"",
		},
		ShowVarsCommand{
			VarLineNum: 8,
			VarName:    "role",
		},
	}
	result, err := DoVarTrace(rego, query, input, nil, commands)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestGetDepTreePretty(t *testing.T) {
	rego := rbacRego
	plan, err := PlanModuleAndGetIr(context.Background(), rego, false, true)
	if err != nil {
		t.Fatal(err)
	}
	treeMap := GetDepTreePretty(plan)
	fmt.Println(treeMap)
}

type staticCallTreeTestData struct {
	description string
	rego        string
	entrypoint  string
	expected    *api.RuleParent
}

func TestGetStaticCallTree(t *testing.T) {
	testData := []staticCallTreeTestData{
		{
			description: "rbac.rego",
			rego:        rbacRego,
			entrypoint:  "allow",
			expected: &api.RuleParent{
				Name:    "allow",
				UID:     "",
				Type:    api.RuleParentTypeParent,
				Default: "false",
				Children: []api.RuleParentChildrenItem{
					{
						Type: api.RuleChildRuleParentChildrenItem,
						RuleChild: api.RuleChild{
							Name:  "allow:30",
							UID:   "",
							Type:  api.RuleChildTypeChild,
							Value: "true",
							Statements: []api.RuleStatement{
								{
									Name: "user_is_admin",
									UID:  "",
									Dependencies: []api.RuleStatementDependenciesItem{
										{
											Type: api.RuleParentRuleStatementDependenciesItem,
											RuleParent: api.RuleParent{
												Name:    "user_is_admin",
												UID:     "",
												Type:    api.RuleParentTypeParent,
												Default: "",
												Children: []api.RuleParentChildrenItem{{
													Type: api.RuleChildRuleParentChildrenItem,
													RuleChild: api.RuleChild{
														Name:  "user_is_admin:43",
														UID:   "",
														Type:  api.RuleChildTypeChild,
														Value: "true",
														Statements: []api.RuleStatement{
															{
																Name: "\"admin\" in data.user_roles[input.user]",
																UID:  "",
																Dependencies: []api.RuleStatementDependenciesItem{
																	{
																		Type:   api.StringRuleStatementDependenciesItem,
																		String: "input.user",
																	}, {
																		Type:   api.StringRuleStatementDependenciesItem,
																		String: "data.user_roles",
																	}},
															},
														},
													},
												}},
											},
										},
									},
								},
							},
						},
					},
					{
						Type: api.RuleChildRuleParentChildrenItem,
						RuleChild: api.RuleChild{
							Name:  "allow:33",
							UID:   "",
							Type:  api.RuleChildTypeChild,
							Value: "true",
							Statements: []api.RuleStatement{
								{
									Name: "some grant in user_is_granted",
									UID:  "",
									Dependencies: []api.RuleStatementDependenciesItem{
										{
											Type: api.RuleParentRuleStatementDependenciesItem,
											RuleParent: api.RuleParent{
												Name:    "user_is_granted",
												UID:     "",
												Type:    api.RuleParentTypeParent,
												Default: "",
												Children: []api.RuleParentChildrenItem{{
													Type: api.RuleChildRuleParentChildrenItem,
													RuleChild: api.RuleChild{
														Name:  "user_is_granted:47",
														UID:   "",
														Type:  api.RuleChildTypeChild,
														Value: "",
														Statements: []api.RuleStatement{
															{
																Name: "some role in data.user_roles[input.user]",
																UID:  "",
																Dependencies: []api.RuleStatementDependenciesItem{
																	{
																		Type:   api.StringRuleStatementDependenciesItem,
																		String: "input.user",
																	},
																	{
																		Type:   api.StringRuleStatementDependenciesItem,
																		String: "data.user_roles",
																	},
																},
															},
															{
																Name: "some grant in data.role_grants[role]",
																UID:  "",
																Dependencies: []api.RuleStatementDependenciesItem{{
																	Type:   api.StringRuleStatementDependenciesItem,
																	String: "data.role_grants",
																}},
															},
														},
													},
												}},
											},
										},
									},
								},
								{
									Name: "input.action == grant.action",
									UID:  "",
									Dependencies: []api.RuleStatementDependenciesItem{
										{
											Type:   api.StringRuleStatementDependenciesItem,
											String: "input.action",
										},
									},
								},
								{
									Name: "input.type == grant.type",
									UID:  "",
									Dependencies: []api.RuleStatementDependenciesItem{
										{
											Type:   api.StringRuleStatementDependenciesItem,
											String: "input.type",
										},
									},
								},
							},
						},
					},
				},
			}},
		{
			description: "deps to nested documents",
			rego: `package test
import future.keywords.if
import data.foo
allow if data.a.b.c[_] = input.a.b.c[foo]`,
			entrypoint: "allow",
			expected: &api.RuleParent{
				Name:    "allow",
				UID:     "",
				Type:    api.RuleParentTypeParent,
				Default: "",
				Children: []api.RuleParentChildrenItem{
					{
						Type: api.RuleChildRuleParentChildrenItem,
						RuleChild: api.RuleChild{
							Name:  "allow:4",
							UID:   "",
							Type:  api.RuleChildTypeChild,
							Value: "true",
							Statements: []api.RuleStatement{
								{
									Name: "data.a.b.c[_] = input.a.b.c[foo]",
									UID:  "",
									Dependencies: []api.RuleStatementDependenciesItem{
										{
											Type:   api.StringRuleStatementDependenciesItem,
											String: "data.foo",
										},
										{
											Type:   api.StringRuleStatementDependenciesItem,
											String: "data.a.b.c",
										},
										{
											Type:   api.StringRuleStatementDependenciesItem,
											String: "input.a.b.c",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			description: "else blocks",
			rego: `package test
import future.keywords.if
allow if {
	2=0
} else if {
	1=0
} else if {
	0=0
}`,
			entrypoint: "allow",
			expected: &api.RuleParent{
				Name:    "allow",
				UID:     "",
				Type:    api.RuleParentTypeParent,
				Default: "",
				Children: []api.RuleParentChildrenItem{
					{
						Type: api.RuleChildElseRuleParentChildrenItem,
						RuleChildElse: api.RuleChildElse{
							Name: "allow",
							UID:  "",
							Type: api.RuleChildElseTypeChildElse,
							Children: []api.RuleChild{
								{
									Name:  "allow:3",
									UID:   "",
									Type:  api.RuleChildTypeChild,
									Value: "true",
									Statements: []api.RuleStatement{
										{
											Name:         "2=0",
											UID:          "",
											Dependencies: nil,
										},
									},
								},
								{
									Name:  "allow:5",
									UID:   "",
									Type:  api.RuleChildTypeChild,
									Value: "true",
									Statements: []api.RuleStatement{
										{
											Name:         "1=0",
											UID:          "",
											Dependencies: nil,
										},
									},
								},
								{
									Name:  "allow:7",
									UID:   "",
									Type:  api.RuleChildTypeChild,
									Value: "true",
									Statements: []api.RuleStatement{
										{
											Name:         "0=0",
											UID:          "",
											Dependencies: nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, data := range testData {
		result, err := GetStaticCallTree(data.rego, data.entrypoint, true)
		if err != nil {
			t.Fatal(data.description, err)
		}

		// Compare the result and expected objects
		if d := cmp.Diff(data.expected, result); d != "" {
			t.Errorf("Result and expected objects differ for \"%s\": (-want +got)\n%s", data.description, d)
		}
	}
}

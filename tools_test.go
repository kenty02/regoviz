package main

import (
	"context"
	_ "embed"
	"fmt"
	"testing"
)

//func TestCompilerPlanTarget(t *testing.T) {
//	files := map[string]string{
//		"test.rego": `# Role-based Access Control (RBAC)
//# --------------------------------
//#
//# This example defines an RBAC model for a Pet Store API. The Pet Store API allows
//# users to look at pets, adopt them, update their stats, and so on. The policy
//# controls which users can perform actions on which resources. The policy implements
//# a classic Role-based Access Control model where users are assigned to roles and
//# roles are granted the ability to perform some action(s) on some type of resource.
//#
//# This example shows how to:
//#
//#	* Define an RBAC model in Rego that interprets role mappings represented in JSON.
//#	* Iterate/search across JSON data structures (e.g., role mappings)
//#
//# For more information see:
//#
//#	* Rego comparison to other systems: https://www.openpolicyagent.org/docs/latest/comparison-to-other-systems/
//#	* Rego Iteration: https://www.openpolicyagent.org/docs/latest/#iteration
//
//package test
//
//import future.keywords.contains
//import future.keywords.if
//import future.keywords.in
//
//# By default, deny requests.
//default allow := false
//
//# Allow admins to do anything.
//allow if user_is_admin
//
//# Allow the action if the user is granted permission to perform the action.
//allow if {
//	# Find grants for the user.
//	some grant in user_is_granted
//
//	# Check if the grant permits the action.
//	input.action == grant.action
//	input.type == grant.type
//}
//
//# user_is_admin is true if "admin" is among the user's roles as per data.user_roles
//user_is_admin if "admin" in data.user_roles[input.user]
//
//# user_is_granted is a set of grants for the user identified in the request.
//# The grant will be contained if the set user_is_granted for every...
//user_is_granted contains grant if {
//	# role assigned an element of the user_roles for this user...
//	some role in data.user_roles[input.user]
//
//	# grant assigned a single grant from the grants list for 'role'...
//	some grant in data.role_grants[role]
//}`,
//	}
//
//	for _, useMemoryFS := range []bool{false, true} {
//		test.WithTestFS(files, useMemoryFS, func(root string, fsys fs.FS) {
//
//			compiler := compile.New().
//				WithFS(fsys).
//				WithPaths(root).
//				WithTarget("plan").
//				WithEntrypoints("test/allow")
//			err := compiler.Build(context.Background())
//			if err != nil {
//				t.Fatal(err)
//			}
//			rs := reflect.ValueOf(compiler).Elem()
//			rsPolicy := rs.FieldByName("policy").Elem()
//			var is *ir.Policy = (*ir.Policy)(unsafe.Pointer(rsPolicy.UnsafeAddr())) // unsafe.Pointer を使うと中身を取れる
//			fmt.Println(is)
//		})
//	}
//
//}

func TestPlan(t *testing.T) {
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
	_, err := plan(ctx, rego, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlanAsText(t *testing.T) {
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

	plan, err := planAsText(ctx, rego, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(plan)
	if len(plan) < 10 {
		t.Fatal("plan is too short")
	}

}

//go:embed samples/rbac.rego
var rbacRego string

func TestGetDepTreeMap(t *testing.T) {
	ctx := context.Background()

	plan, err := plan(ctx, rbacRego, false)
	if err != nil {
		t.Fatal(err)
	}

	treeMap := getDepTreeMap(plan)

	fmt.Println(treeMap)

	if len(treeMap) < 1 {
		t.Fatal("treeMap is too short")
	}
}

func TestGetMermaidFlowchart(t *testing.T) {
	ctx := context.Background()

	plan, err := plan(ctx, rbacRego, false)
	if err != nil {
		t.Fatal(err)
	}

	mermaid := getMermaidFlowchart(plan)

	url, err := getMermaidUrl(mermaid, true)
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

func TestRegoVarTrace(t *testing.T) {
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
		//FixVarCommand{
		//	varLineNum: 8,
		//	varName:    "role",
		//	varValue:   "\"hoge\"",
		//},
		ShowVarsCommand{
			varLineNum: 8,
			varName:    "role",
		},
	}
	result, err := regoVarTrace(rego, query, input, nil, commands)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestGetDepTreePretty(t *testing.T) {
	rego := rbacRego
	plan, err := plan(context.Background(), rego, false)
	if err != nil {
		t.Fatal(err)
	}
	treeMap := getDepTreePretty(plan)
	fmt.Println(treeMap)
}

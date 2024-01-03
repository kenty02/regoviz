package astprint

import (
	_ "embed"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"go.uber.org/zap/buffer"
	"testing"
)

//go:embed testdata/rbac.rego
var rbacRego string

func TestPretty(t *testing.T) {
	const moduleName = "my_module"
	// Parse the input module to obtain the AST representation.
	mod, err := ast.ParseModule(moduleName, rbacRego)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new compiler instance and compile the module.
	c := ast.NewCompiler()

	mods := map[string]*ast.Module{
		moduleName: mod,
	}

	if c.Compile(mods); c.Failed() {
		t.Fatal(c.Errors)
	}

	buf := buffer.Buffer{}
	err = Pretty(&buf, mod)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
}

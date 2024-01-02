package astprint

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
)

// Visitor defines the interface for visiting IR nodes.
type Visitor interface {
	Before(x interface{})
	Visit(x interface{}, stringer func(x interface{}) string) (Visitor, error)
	After(x interface{})
}

// Walk invokes the visitor for nodes under x.
func Walk(vis Visitor, x interface{}) error {
	impl := walkerImpl{
		vis: vis,
	}
	impl.walk(x, nil)
	return impl.err
}

type walkerImpl struct {
	vis Visitor
	err error
}

func (w *walkerImpl) walk(x interface{}, stringer func(x interface{}) string) {
	if w.err != nil { // abort on error
		return
	}
	if x == nil {
		return
	}

	prev := w.vis
	w.vis.Before(x)
	defer w.vis.After(x)
	w.vis, w.err = w.vis.Visit(x, stringer)
	if w.err != nil {
		return
	} else if w.vis == nil {
		w.vis = prev
		return
	}

	switch x := x.(type) {
	case *ast.Module:
		w.walk(x.Package, nil)
		w.walk(x.Rules, nil)
		w.walk(x.Comments, nil)
	case *ast.Package:
		w.walk(x.Path, nil)
		w.walk(x.Location, nil)
	case []*ast.Rule:
		for _, r := range x {
			w.walk(r, nil)
		}
	case *ast.Rule:
		w.walk(x.Default, func(x interface{}) string {
			y := x.(bool)
			return fmt.Sprintf("Is default rule: %v", y)
		})
		w.walk(x.Head, nil)
		w.walk(x.Body, nil)
		if x.Else != nil {
			w.walk(x.Else, nil)
		}
	case ast.Body:
		for _, e := range x {
			w.walk(e, nil)
		}
	case *ast.Expr:
		w.walk(x.Terms, nil)
	case []*ast.Term:
		for _, t := range x {
			w.walk(t, nil)
		}
	case []*ast.Comment:
		for _, c := range x {
			w.walk(c, func(x interface{}) string {
				y := x.(*ast.Comment)
				return fmt.Sprintf("%T (%d,%d) \"%s\"", y, y.Location.Row, y.Location.Col, string(y.Text))
			})
		}
		//case *ast.Static:
		//	for _, s := range x.Strings {
		//		w.walk(s)
		//	}
		//	for _, f := range x.BuiltinFuncs {
		//		w.walk(f)
		//	}
		//	for _, f := range x.Files {
		//		w.walk(f)
		//	}
		//case *Plans:
		//	for _, pl := range x.Plans {
		//		w.walk(pl)
		//	}
		//case *Funcs:
		//	for _, fn := range x.Funcs {
		//		w.walk(fn)
		//	}
		//case *Func:
		//	for _, b := range x.Blocks {
		//		w.walk(b)
		//	}
		//case *Plan:
		//	for _, b := range x.Blocks {
		//		w.walk(b)
		//	}
		//case *Block:
		//	for _, s := range x.Stmts {
		//		w.walk(s)
		//	}
		//case *BlockStmt:
		//	for _, b := range x.Blocks {
		//		w.walk(b)
		//	}
		//case *ScanStmt:
		//	w.walk(x.Block)
		//case *NotStmt:
		//	w.walk(x.Block)
		//case *WithStmt:
		//	w.walk(x.Block)
		//}

	}
}

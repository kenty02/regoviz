package directDeps

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/dependencies"
	"sort"
)

// similar to opa/dependencies but only returns direct dependencies

// Copyright 2017 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

// Base returns the list of base data documents that the given AST element depends on.
//
// The returned refs are always constant and are truncated at any point where they become
// dynamic. That is, a ref like data.a.b[x] will be truncated to data.a.b.
func Base(compiler *ast.Compiler, x interface{}) ([]ast.Ref, error) {
	baseRefs, err := base(compiler, x)
	if err != nil {
		return nil, err
	}

	return dedup(baseRefs), nil
}

func base(compiler *ast.Compiler, x interface{}) ([]ast.Ref, error) {
	refs, err := dependencies.Minimal(x)
	if err != nil {
		return nil, err
	}

	var baseRefs []ast.Ref
	for _, r := range refs {
		r = r.ConstantPrefix()
		if rules := compiler.GetRules(r); len(rules) > 0 {
			//for _, rule := range rules {
			//	bases, err := base(compiler, rule)
			//	if err != nil {
			//		panic("not reached")
			//	}
			//
			//	baseRefs = append(baseRefs, bases...)
			//}
		} else {
			baseRefs = append(baseRefs, r)
		}
	}

	return baseRefs, nil
}

// Virtual returns the list of virtual data documents that the given AST element depends
// on.
//
// The returned refs are always constant and are truncated at any point where they become
// dynamic. That is, a ref like data.a.b[x] will be truncated to data.a.b.
func Virtual(compiler *ast.Compiler, x interface{}) ([]ast.Ref, error) {
	virtualRefs, err := virtual(compiler, x)
	if err != nil {
		return nil, err
	}

	return dedup(virtualRefs), nil
}

func virtual(compiler *ast.Compiler, x interface{}) ([]ast.Ref, error) {
	refs, err := dependencies.Minimal(x)
	if err != nil {
		return nil, err
	}

	var virtualRefs []ast.Ref
	for _, r := range refs {
		r = r.ConstantPrefix()
		if rules := compiler.GetRules(r); len(rules) > 0 {
			for _, rule := range rules {
				//virtuals, err := virtual(compiler, rule)
				//if err != nil {
				//	panic("not reached")
				//}

				virtualRefs = append(virtualRefs, rule.Path())
				//virtualRefs = append(virtualRefs, virtuals...)
			}
		}
	}

	return virtualRefs, nil
}

func dedup(refs []ast.Ref) []ast.Ref {
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Compare(refs[j]) < 0
	})

	return filter(refs, func(a, b ast.Ref) bool {
		return a.Compare(b) == 0
	})
}

// filter removes all items from the list that cause pref to return true. It is
// called on adjacent pairs of elements, and the one passed as the second argument
// to pref is considered the current one being examined. The first argument will
// be the element immediately preceding it.
func filter(rs []ast.Ref, pred func(ast.Ref, ast.Ref) bool) (filtered []ast.Ref) {
	if len(rs) == 0 {
		return nil
	}

	last := rs[0]
	filtered = append(filtered, last)
	for i := 1; i < len(rs); i++ {
		cur := rs[i]
		if pred(last, cur) {
			continue
		}

		filtered = append(filtered, cur)
		last = cur
	}

	return filtered
}

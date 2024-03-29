package analyzer

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/compile"
	"github.com/open-policy-agent/opa/ir"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/topdown/print"
	treemap "github.com/xlab/treeprint"
	"math/rand"
	"reflect"
	"regexp"
	"regoviz/internal/analyzer/directDeps"
	"regoviz/internal/api"
	"regoviz/internal/utils"
	"slices"
	"strings"
)

func CompileModuleStringToAst(moduleCode string, withPrintStatements bool, withStrict bool) (*ast.Module, *ast.Compiler, error) {
	const moduleName = "my_module"
	// Parse the input module to obtain the AST representation.
	mod, err := ast.ParseModule(moduleName, moduleCode)
	if err != nil {
		return nil, nil, err
	}

	// Create a new compiler instance and compile the module.
	c := ast.NewCompiler().WithEnablePrintStatements(withPrintStatements).WithStrict(withStrict)

	mods := map[string]*ast.Module{
		moduleName: mod,
	}

	if c.Compile(mods); c.Failed() {
		return nil, nil, c.Errors
	}

	return c.Modules[moduleName], c, nil
}

func PlanModuleAndGetIr(ctx context.Context, rego string, print bool, autoDetermineEntrypoints bool) (*ir.Policy, error) {
	mod, err := ast.ParseModuleWithOpts("a.rego", rego, ast.ParserOptions{ProcessAnnotation: true})
	if err != nil {
		return nil, err
	}
	b := &bundle.Bundle{
		Modules: []bundle.ModuleFile{
			{
				URL:    "/url",
				Path:   "/a.rego",
				Raw:    []byte(rego),
				Parsed: mod,
			},
		},
	}

	compiler := compile.New().
		WithTarget(compile.TargetPlan).
		WithBundle(b).
		WithRegoAnnotationEntrypoints(true).
		WithEnablePrintStatements(print)
	if autoDetermineEntrypoints {
		moduleAst, _, err := CompileModuleStringToAst(rego, false, true)

		if err != nil {
			return nil, err
		}
		var entrypoints []string
		for _, rule := range moduleAst.Rules {
			pkg := moduleAst.Package
			path := make(ast.Ref, len(pkg.Path)-1)
			path[0] = ast.VarTerm(string(pkg.Path[1].Value.(ast.String)))
			copy(path[1:], pkg.Path[2:])
			pathString := path.String()
			// replace . to /
			pathString = strings.ReplaceAll(pathString, ".", "/")
			entrypoint := fmt.Sprintf("%s/%s", pathString, rule.Head.Name.String())
			entrypoints = append(entrypoints, entrypoint)
		}
		compiler = compiler.WithEntrypoints(entrypoints...)
	}
	if err := compiler.Build(ctx); err != nil {
		return nil, err
	}

	compiledBundle := compiler.Bundle()
	var policy ir.Policy

	if err := json.Unmarshal(compiledBundle.PlanModules[0].Raw, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

func PlanAsText(ctx context.Context, rego string, print bool) (string, error) {
	policy, err := PlanModuleAndGetIr(ctx, rego, print, true)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	if err := ir.Pretty(&buf, policy); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetDepTreeMap(policy *ir.Policy) map[string][]string {
	depTreeMap := map[string][]string{}
	for _, fun := range policy.Funcs.Funcs {
		caller := fun.Name
		for _, block := range fun.Blocks {
			for _, stmt := range block.Stmts {
				switch stmt := stmt.(type) {
				case *ir.CallStmt:
					callee := stmt.Func
					depTreeMap[caller] = append(depTreeMap[caller], callee)
				}
			}
		}
	}
	return depTreeMap
}

func GetDepTreePretty(policy *ir.Policy) string {
	depTreeMap := GetDepTreeMap(policy)
	tree := treemap.New()
	for caller, callees := range depTreeMap {
		callerNode := tree.AddBranch(caller)
		for _, callee := range callees {
			callerNode.AddNode(callee)
		}
	}
	return tree.String()
}

func esc(s any) string {
	str := fmt.Sprintf("%s", s)
	str = strings.ReplaceAll(str, "\"", "#quot;")
	return str
}
func GetMermaidFlowchart(policy *ir.Policy) string {
	buf := bytes.Buffer{}
	indentLevel := 0
	nodeIdCounter := 0
	subgraphIdCounter := 0
	var add = func(s string) {

		for i := 0; i < indentLevel; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString(s)
		buf.WriteString("\n")
	}
	add("flowchart TB")
	indentLevel++
	for _, fun := range policy.Funcs.Funcs {
		add(fmt.Sprintf(`subgraph "%s(params: TODO)"`, fun.Name))
		indentLevel++
		var addBlocks func(blocks []*ir.Block)
		addBlocks = func(blocks []*ir.Block) {
			for blockIdx, block := range blocks {
				add(fmt.Sprintf(`subgraph "Block %d/%d"`, blockIdx+1, len(blocks)))
				indentLevel++
				firstStmt := true
				var getStmtNodeId = func() int {
					id := nodeIdCounter
					nodeIdCounter++
					return id
				}
				var connectNodes = func(from, to int) {
					add(fmt.Sprintf(`node%d --> node%d`, from, to))
				}
				var connectLastTwoNodes = func() {
					if !firstStmt {
						connectNodes(nodeIdCounter-2, nodeIdCounter-1)
					}
				}
				for _, stmt := range block.Stmts {
					switch stmt := stmt.(type) {
					case *ir.CallStmt:
						var sb strings.Builder
						for i, arg := range stmt.Args {
							if i > 0 {
								sb.WriteString(", ")
							}
							sb.WriteString(arg.Value.String())
						}
						stmtNodeId := getStmtNodeId()
						add(fmt.Sprintf(`node%d["%s(%s)"]`, stmtNodeId, stmt.Func, sb.String()))
						connectLastTwoNodes()
					case *ir.ScanStmt:
						stmtNodeId := getStmtNodeId()
						add(fmt.Sprintf(`node%d["for (%s, %s) in %s {TODO}"]`, stmtNodeId, stmt.Key.String(), stmt.Value.String(), stmt.Source.String()))
						add(fmt.Sprintf(`subgraph "ScanStmt%d"`, subgraphIdCounter))
						subgraphIdCounter++
						prevNodeId := nodeIdCounter - 1
						blockStmtStartNodeId := nodeIdCounter
						indentLevel++
						// create array of blocks to pass to addBlocks
						blocks := make([]*ir.Block, 1)
						blocks[0] = stmt.Block
						addBlocks(blocks)
						indentLevel--
						if nodeIdCounter > blockStmtStartNodeId {
							connectNodes(prevNodeId, blockStmtStartNodeId)
						}
						add("end")
					case *ir.BlockStmt:
						add(fmt.Sprintf(`subgraph "BlockStmt%d"`, subgraphIdCounter))
						subgraphIdCounter++
						prevNodeId := nodeIdCounter - 1
						blockStmtStartNodeId := nodeIdCounter
						indentLevel++
						addBlocks(stmt.Blocks)
						indentLevel--
						if nodeIdCounter > blockStmtStartNodeId {
							connectNodes(prevNodeId, blockStmtStartNodeId)
						}
						add("end")
					default:
						stmtName := reflect.TypeOf(stmt).String()
						//marshal stmt
						stmtBytes, err := json.Marshal(stmt)
						if err != nil {
							panic(err)
						}
						// if it starts with "*ir.", remove it. if not, warn
						if stmtName[:4] == "*ir." {
							stmtName = "TODO(" + stmtName[4:] + ")"
						} else {
							stmtName = "TODO(warn: not a stmt!): " + stmtName
						}
						// write stmt type name as is
						stmtNodeId := getStmtNodeId()
						add(fmt.Sprintf(`node%d["%s \n%s"]`, stmtNodeId, stmtName, esc(stmtBytes)))
						connectLastTwoNodes()
					}
					if firstStmt {
						firstStmt = false
					}
				}
				indentLevel--
				add("end")
			}
		}
		addBlocks(fun.Blocks)
		indentLevel--
		add("end")
	}
	indentLevel--
	return buf.String()
}

type MermaidConfig struct {
	Theme string `json:"theme"`
}
type MermaidPan struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type MermaidState struct {
	Code          string        `json:"code"`
	Mermaid       MermaidConfig `json:"mermaid"`
	AutoSync      bool          `json:"autoSync"`
	UpdateDiagram bool          `json:"updateDiagram"`
	Pan           MermaidPan    `json:"pan"`
	Zoom          float64       `json:"zoom"`
	PanZoom       bool          `json:"panZoom"`
}

func GetMermaidUrl(code string, edit bool) (string, error) {
	state := MermaidState{
		Code: code,
		Mermaid: MermaidConfig{
			Theme: "default",
		},
		AutoSync:      true,
		UpdateDiagram: true,
		Pan: MermaidPan{
			X: -71.84227712898122,
			Y: 196.46405886972394,
		},
		Zoom:    5.001651287078857,
		PanZoom: true,
	}
	// encode to json
	st, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	// compress code
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err = w.Write(st)
	if err != nil {
		return "", nil
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	// encode to base64
	sEnc := base64.StdEncoding.EncodeToString(b.Bytes())
	// create url
	var baseUrl string
	if edit {
		baseUrl = "https://mermaid.live/edit#pako:"
	} else {
		baseUrl = "https://mermaid.live/view#pako:"
	}
	return baseUrl + sEnc, nil
}

func evalRegoWithPrint(code, query string, input, data map[string]interface{}) (*rego.ResultSet, string, error) {
	var buf bytes.Buffer

	maybeInputFunc := func(r *rego.Rego) {}
	maybeStoreFunc := func(r *rego.Rego) {}

	if input != nil {
		maybeInputFunc = rego.Input(input)
	}
	if data != nil {
		// Manually create the storage layer. inmem.NewFromObject returns an
		// in-memory store containing the supplied data.
		store := inmem.NewFromObject(data)
		maybeStoreFunc = rego.Store(store)
	}
	r := rego.New(
		rego.Query(query),
		rego.Module("example.rego", code),
		maybeInputFunc,
		rego.EnablePrintStatements(true),
		rego.PrintHook(topdown.NewPrintHook(&buf)),
		maybeStoreFunc,
		rego.UnsafeBuiltins(map[string]struct{}{"http.send": {}}),
	)

	rs, err := r.Eval(context.Background())
	if err != nil {
		return nil, "", err
	}

	return &rs, buf.String(), nil
}

func EvalRegoWithPrintAndTrace(code, query string, input, data map[string]interface{}, p print.Hook, t topdown.QueryTracer) (*rego.ResultSet, error) {

	maybeInputFunc := func(r *rego.Rego) {}
	maybeStoreFunc := func(r *rego.Rego) {}

	if input != nil {
		maybeInputFunc = rego.Input(input)
	}
	if data != nil {
		// Manually create the storage layer. inmem.NewFromObject returns an
		// in-memory store containing the supplied data.
		store := inmem.NewFromObject(data)
		maybeStoreFunc = rego.Store(store)
	}
	r := rego.New(
		rego.Query(query),
		rego.Module("example.rego", code),
		maybeInputFunc,
		rego.EnablePrintStatements(true),
		rego.PrintHook(p),
		rego.QueryTracer(t),
		rego.Trace(true),
		maybeStoreFunc,
		rego.UnsafeBuiltins(map[string]struct{}{"http.send": {}}),
	)

	rs, err := r.Eval(context.Background())
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

type CodeInject struct {
	line    int
	code    string
	replace bool
}

func injectCode(code string, cis []CodeInject) string {
	// コードを行に分割する
	lines := strings.Split(code, "\n")

	// 各挿入について処理
	for _, ci := range cis {
		if strings.Contains(ci.code, "\n") || strings.Contains(ci.code, "\r") {
			panic("ci.code contains new line")
		}
		if ci.line > 0 && ci.line <= len(lines) {
			// 行番号は1から始まるが、スライスのインデックスは0から始まるため調整
			index := ci.line - 1
			if ci.replace {
				// 指定された行を置換
				lines[index] = ci.code
			} else {
				// 指定された行にコードを挿入
				lines[index] = lines[index] + ci.code
			}
		}
	}

	// 更新されたコード行を再結合
	return strings.Join(lines, "\n")
}

func replaceTokenInLine(line, token, replacement string) (string, error) {
	// トークンの正規表現パターンを作成します。
	// 単語の境界 (\b) を使用して、完全なトークンのみをマッチさせます。
	tokenPattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(token))
	re, err := regexp.Compile(tokenPattern)
	if err != nil {
		return "", err
	}

	// 置換を実行します。
	return re.ReplaceAllString(line, replacement), nil
}

type ShowVarsCommand struct {
	VarLineNum int
	VarName    string
}

type FixVarCommand struct {
	VarLineNum int
	VarName    string
	VarValue   string
}

// varトークンにマッチする正規表現パターン
// \bは単語の境界を意味し、hogeが単独のトークンとして現れる場合にのみマッチする
var varPattern = `\b[A-Za-z_][A-Za-z_0-9]*\b`
var varRe = regexp.MustCompile(varPattern)

// 正規表現パターン: VTBEGINとVTENDの間にあるvarNameとvarValueを抽出
var vtPattern = `VTBEGIN ([A-Za-z_][A-Za-z_0-9]*) (.+) VTEND`
var vtRe = regexp.MustCompile(vtPattern)

func DoVarTrace(code, query string, input, data map[string]interface{}, commands []interface{}) (string, error) {
	lines := strings.Split(strings.ReplaceAll(code, "\r\n", "\n"), "\n")

	var sb strings.Builder
	var showVarsCommands []ShowVarsCommand
	var fixVarCommands []FixVarCommand

	for _, command := range commands {
		switch command := command.(type) {
		case ShowVarsCommand:
			showVarsCommands = append(showVarsCommands, command)
		case FixVarCommand:
			fixVarCommands = append(fixVarCommands, command)
		default:
			return "", fmt.Errorf("unknown command type: %T", command)
		}
	}
	modifiedCode := code

	// iterate over commands
	for _, command := range fixVarCommands {
		// 行番号は1から始まるが、スライスのインデックスは0から始まるため調整
		index := command.VarLineNum - 1
		line := lines[index]

		// マッチがあるかどうか調べる
		matches := varRe.FindAllString(line, -1)
		// check if varName is in matches
		varNameFound := false
		for _, match := range matches {
			if match == command.VarName {
				varNameFound = true
				// random suffixed var name
				suffix := rand.Intn(100000)
				newVarName := fmt.Sprintf("%s_fixed_%d", command.VarName, suffix)
				// そのマッチを置換
				var err error
				line, err = replaceTokenInLine(line, command.VarName, newVarName) // TODO: BNFと一貫性のある置換
				if err != nil {
					return "", err
				}
			}
		}
		if !varNameFound {
			return "", fmt.Errorf("var %s not found in line %d", command.VarName, command.VarLineNum)
		}
		line = line + fmt.Sprintf("; %s:=%s", command.VarName, command.VarValue)
		cis := []CodeInject{
			{command.VarLineNum, line, true}}
		modifiedCode = injectCode(modifiedCode, cis)
		// just test if it works
		_, _, err := evalRegoWithPrint(modifiedCode, query, input, nil)
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprintf("変数%sの値を%sに固定しました。\n", command.VarName, command.VarValue))
	}
	for _, command := range showVarsCommands {
		// 行番号は1から始まるが、スライスのインデックスは0から始まるため調整
		index := command.VarLineNum - 1
		line := lines[index]

		// マッチがあるかどうか調べる
		matches := varRe.FindAllString(line, -1)
		// check if varName is in matches
		varNameFound := false
		for _, match := range matches {
			if match == command.VarName {
				varNameFound = true
				break
			}
		}
		if !varNameFound {
			return "", fmt.Errorf("var %s not found in line %d", command.VarName, command.VarLineNum)
		}

		cis := []CodeInject{
			{command.VarLineNum, fmt.Sprintf(";print(\"VTBEGIN %s\", %s, \"VTEND\");false", command.VarName, command.VarName), false},
		}
		injectedCode := injectCode(modifiedCode, cis)
		_, printed, err := evalRegoWithPrint(injectedCode, query, input, data)
		if err != nil {
			return "", err
		}

		// 文字列がパターンにマッチするかどうかをチェック
		vtMatches := vtRe.FindAllStringSubmatch(printed, -1)

		// varName to set of varValue map
		varSet := mapset.NewSet[string]()
		// マッチした結果からvarNameとvarValueを抽出して表示
		for _, match := range vtMatches {
			if len(match) > 2 {
				varName := match[1]
				varValue := match[2]

				if varName != command.VarName {
					return "", fmt.Errorf("expected varName: %s, actual varName: %s", command.VarName, varName)
				}
				varSet.Add(varValue)
			}
		}

		cis2 := []CodeInject{
			{command.VarLineNum, fmt.Sprintf(";print(\"VTBEGIN %s\", %s, \"VTEND\")", command.VarName, command.VarName), false},
		}
		injectedCode = injectCode(modifiedCode, cis2)
		_, printed, err = evalRegoWithPrint(injectedCode, query, input, nil)
		if err != nil {
			return "", err
		}

		actualVarSet := mapset.NewSet[string]()
		vtMatches2 := vtRe.FindAllStringSubmatch(printed, -1)
		for _, match := range vtMatches2 {
			if len(match) > 2 {
				varName := match[1]
				varValue := match[2]

				if varName != command.VarName {
					return "", fmt.Errorf("expected varName: %s, actual varName: %s", command.VarName, varName)
				}
				actualVarSet.Add(varValue)
			}
		}

		if !varSet.IsEmpty() {
			sb.WriteString(fmt.Sprintf("変数%sは%sの値を取り得ます。", command.VarName, varSet.ToSlice()))
			if !actualVarSet.IsEmpty() {
				sb.WriteString(fmt.Sprintf("評価パスで実際に代入された値の集合は%sです。", actualVarSet.ToSlice()))
			}
		} else {
			sb.WriteString(fmt.Sprintf("変数%sは値を取り得ませんでした。", command.VarName))
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

type UIDType int

const (
	UIDTypeEmpty UIDType = iota
	UIDTypeRandom
	UIDTypeRandomWithLocation
	UIDTypeDeterministicWithLocation
)

func GetAvailableEntrypointsForCallTree(code string) ([]string, error) {
	moduleAst, _, err := CompileModuleStringToAst(code, false, false)
	if err != nil {
		return nil, err
	}
	entrypointsSet := map[string]struct{}{}
	for _, rule := range moduleAst.Rules {
		entrypointsSet[rule.Head.Name.String()] = struct{}{}
	}
	var entrypoints []string
	for entrypoint := range entrypointsSet {
		entrypoints = append(entrypoints, entrypoint)
	}
	slices.SortFunc(entrypoints, func(a, b string) int {
		return strings.Compare(a, b)
	})

	return entrypoints, nil
}

func GetStaticCallTree(code, entrypoint string, uidType UIDType) (*api.RuleParent, []interface{}, error) {
	moduleAst, moduleAstCompiler, err := CompileModuleStringToAst(code, true, true)
	if err != nil {
		return nil, nil, err
	}

	uidCounter := 0
	maybeGenerateUID := func(locationOrNil *ast.Location) string {
		if uidType == UIDTypeEmpty {
			return ""
		} else if uidType == UIDTypeRandom {
			return utils.GenerateId()
		} else if uidType == UIDTypeRandomWithLocation {
			locationString := "nil"
			if locationOrNil != nil {
				locationString = locationOrNil.String()
			}
			return fmt.Sprintf("%s_%s", locationString, utils.GenerateId())
		} else if uidType == UIDTypeDeterministicWithLocation {
			locationString := "nil"
			if locationOrNil != nil {
				locationString = locationOrNil.String()
			}
			uid := fmt.Sprintf("%s_%d", locationString, uidCounter)
			uidCounter++
			return uid
		} else {
			panic("unknown uidType")
		}
	}

	var nodes []interface{}
	var getRuleParent func(ruleRef ast.Ref) (*api.RuleParent, error)
	getRuleName := func(rule *ast.Rule) string {
		return rule.Head.Name.String()
	}
	// always returns RuleChild even if "rule.Else" is present
	getRuleChild := func(rule *ast.Rule) (*api.RuleChild, error) {
		var statements []api.RuleStatement
		// expecting generated expr(s) to be prepended
		var prependedExprs []*ast.Expr
		for _, expr := range rule.Body {
			if expr.Generated {
				prependedExprs = append(prependedExprs, expr)
				continue
			}
			var deps []api.RuleStatementDependenciesItem
			for _, expr := range append(prependedExprs, expr) {
				base, err := directDeps.Base(moduleAstCompiler, expr)
				if err != nil {
					return nil, err
				}
				for _, ref := range base {
					if ref.HasPrefix(ast.InputRootRef) || ref.HasPrefix(ast.DefaultRootRef) {
						deps = append(deps, api.RuleStatementDependenciesItem{
							Type:   api.StringRuleStatementDependenciesItem,
							String: ref.String(),
						})
					}
				}
				virtual, err := directDeps.Virtual(moduleAstCompiler, expr)
				if err != nil {
					return nil, err
				}
				for _, ref := range virtual {
					ruleParent, err := getRuleParent(ref)
					if err != nil {
						return nil, err
					}
					deps = append(deps, api.RuleStatementDependenciesItem{
						Type:       api.RuleParentRuleStatementDependenciesItem,
						RuleParent: *ruleParent,
					})
				}
			}
			prependedExprs = nil

			statement := api.RuleStatement{
				Name:         string(expr.Location.Text),
				UID:          maybeGenerateUID(expr.Location),
				Dependencies: deps,
				Location: api.NewOptNodeLocation(api.NodeLocation{
					Row: expr.Location.Row,
					Col: expr.Location.Col,
				}),
			}
			statements = append(statements, statement)
			nodes = append(nodes, &statement)

		}
		value := ""
		if rule.Head.Value != nil {
			value = rule.Head.Value.String()
		}
		result := &api.RuleChild{
			Name:       fmt.Sprintf("%s:%d", getRuleName(rule), rule.Head.Location.Row),
			UID:        maybeGenerateUID(rule.Head.Location),
			Type:       api.RuleChildTypeChild,
			Value:      value,
			Statements: statements,
			Location: api.NewOptNodeLocation(api.NodeLocation{
				Row: rule.Head.Location.Row,
				Col: rule.Head.Location.Col,
			}),
		}
		nodes = append(nodes, result)
		return result, nil
	}

	// must pass non-default rules
	getRuleParentChild := func(rule *ast.Rule) (*api.RuleParentChildrenItem, error) {
		ruleChild, err := getRuleChild(rule)
		if err != nil {
			return nil, err
		}

		ruleChildElseChildren := []api.RuleChild{*ruleChild}
		current := rule.Else
		for current != nil {
			ruleChild, err := getRuleChild(current)
			if err != nil {
				return nil, err
			}
			ruleChildElseChildren = append(ruleChildElseChildren, *ruleChild)
			current = current.Else
		}
		var result *api.RuleParentChildrenItem
		if len(ruleChildElseChildren) > 1 {
			result = &api.RuleParentChildrenItem{
				Type: api.RuleChildElseRuleParentChildrenItem,
				RuleChildElse: api.RuleChildElse{
					Name:     getRuleName(rule),
					UID:      maybeGenerateUID(nil),
					Type:     api.RuleChildElseTypeChildElse,
					Children: ruleChildElseChildren,
					Location: api.NewOptNodeLocation(api.NodeLocation{
						Row: ruleChild.Location.Value.Row,
						Col: ruleChild.Location.Value.Col,
					}),
				},
			}
		} else {
			result = &api.RuleParentChildrenItem{
				Type:      api.RuleChildRuleParentChildrenItem,
				RuleChild: *ruleChild,
			}
		}
		nodes = append(nodes, result)
		return result, nil
	}
	getRuleParentChildren := func(rules []*ast.Rule) ([]api.RuleParentChildrenItem, error) {
		var results []api.RuleParentChildrenItem
		for _, rule := range rules {
			child, err := getRuleParentChild(rule)
			if err != nil {
				return nil, err
			}
			results = append(results, *child)
		}
		return results, nil
	}
	getRuleParent = func(ruleRef ast.Ref) (*api.RuleParent, error) {
		var nonDefaultRules []*ast.Rule
		var defaultRule *ast.Rule
		for _, rule := range moduleAstCompiler.GetRulesExact(ruleRef) {
			if rule.Default {
				if defaultRule != nil {
					return nil, fmt.Errorf("multiple default rules found: %s", ruleRef)
				}
				defaultRule = rule
			} else {
				nonDefaultRules = append(nonDefaultRules, rule)
			}
		}
		if len(nonDefaultRules) == 0 {
			return nil, fmt.Errorf("rule not found: %s", ruleRef)
		}
		defaultRuleValue := ""
		if defaultRule != nil {
			defaultRuleValue = defaultRule.Head.Value.String()
		}
		children, err := getRuleParentChildren(nonDefaultRules)
		if err != nil {
			return nil, err
		}
		result := &api.RuleParent{
			Name:     getRuleName(nonDefaultRules[0]),
			UID:      maybeGenerateUID(nil),
			Type:     api.RuleParentTypeParent,
			Default:  defaultRuleValue,
			Children: children,
			Ref:      ruleRef.String(),
		}
		nodes = append(nodes, result)
		return result, nil
	}

	// get entrypoint rule
	var entrypointRuleRef ast.Ref
	for _, rule := range moduleAst.Rules {
		// todo allow absolute name like "data.example.allow"
		if rule.Head.Name.String() == entrypoint {
			entrypointRuleRef = rule.Ref()
			break
		}
	}
	if entrypointRuleRef == nil {
		return nil, nil, fmt.Errorf("entrypoint rule not found: %s", entrypoint)
	}
	result, err := getRuleParent(entrypointRuleRef)
	if err != nil {
		return nil, nil, err
	}

	return result, nodes, nil
}

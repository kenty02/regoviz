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
	treemap "github.com/xlab/treeprint"
	"math/rand"
	"reflect"
	"regexp"
	"regoviz/internal/api"
	"regoviz/internal/utils"
	"strings"
)

func CompileModuleStringToAst(moduleCode string) (*ast.Module, error) {
	const moduleName = "my_module"
	// Parse the input module to obtain the AST representation.
	mod, err := ast.ParseModule(moduleName, moduleCode)
	if err != nil {
		return nil, err
	}

	// Create a new compiler instance and compile the module.
	c := ast.NewCompiler()

	mods := map[string]*ast.Module{
		moduleName: mod,
	}

	if c.Compile(mods); c.Failed() {
		return nil, c.Errors
	}
	//
	//fmt.Println("Expr 1:", c.Modules[moduleName].Rules[0].Body[0])
	//fmt.Println("Expr 2:", c.Modules[moduleName].Rules[0].Body[1])
	//fmt.Println("Expr 3:", c.Modules[moduleName].Rules[0].Body[2])
	//fmt.Println("Expr 4:", c.Modules[moduleName].Rules[0].Body[3])

	return c.Modules[moduleName], nil
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
		moduleAst, err := CompileModuleStringToAst(rego)

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
	)

	rs, err := r.Eval(context.Background())
	if err != nil {
		return nil, "", err
	}

	return &rs, buf.String(), nil
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

func flattenIrStmts(stmts []ir.Stmt) []ir.Stmt {
	var result []ir.Stmt
	for _, stmt := range stmts {
		result = append(result, stmt)
		switch stmt := stmt.(type) {
		case *ir.BlockStmt:
			for _, block := range stmt.Blocks {
				result = append(result, flattenIrStmts(block.Stmts)...)
			}
		case *ir.ScanStmt:
			result = append(result, flattenIrStmts(stmt.Block.Stmts)...)
		case *ir.NotStmt:
			result = append(result, flattenIrStmts(stmt.Block.Stmts)...)
		case *ir.WithStmt:
			result = append(result, flattenIrStmts(stmt.Block.Stmts)...)
		default:
		}
	}
	return result
}

func GetStaticCallTree(code, entrypoint string, useEmptyUID bool) (*api.RuleParent, error) {
	moduleAst, err := CompileModuleStringToAst(code)
	if err != nil {
		return nil, err
	}
	moduleIr, err := PlanModuleAndGetIr(context.Background(), code, false, true)
	if err != nil {
		return nil, err
	}
	//// get entrypoint rule
	//var entrypointRules []*ast.Rule
	//for _, rule := range moduleAst.Rules {
	//	if rule.Head.Name.String() == entrypoint {
	//		entrypointRules = append(entrypointRules, rule)
	//	}
	//}

	ruleNameToFunc := map[string]*ir.Func{}
	funcNameToRuleName := map[string]string{}
	for _, plan := range moduleIr.Plans.Plans {
		packageName := moduleAst.Package.Path.String()
		// if contains [, something is wrong
		if strings.Contains(packageName, "[") {
			return nil, fmt.Errorf("moduleAst.Package.Path.String() contains [")
		}
		ruleName := strings.TrimPrefix(strings.ReplaceAll(plan.Name, "/", "."), strings.TrimPrefix(packageName+".", "data."))
		if len(plan.Blocks) != 1 {
			return nil, fmt.Errorf("len(plan.Blocks) != 1")
		}
		block := plan.Blocks[0]

		var callees []*ir.Func
		for _, stmt := range block.Stmts {
			switch stmt := stmt.(type) {
			case *ir.CallStmt:
				funcName := stmt.Func
				var callee *ir.Func
				for _, fun := range moduleIr.Funcs.Funcs {
					if fun.Name == funcName {
						callee = fun
						break
					}
				}
				if callee == nil {
					return nil, fmt.Errorf("callee not found: %s", funcName)
				}
				callees = append(callees, callee)
			}
		}
		if len(callees) != 1 {
			return nil, fmt.Errorf("len(callees) != 1")
		}
		ruleNameToFunc[ruleName] = callees[0]
		funcNameToRuleName[callees[0].Name] = ruleName
	}

	maybeGenerateUID := func() string {
		if useEmptyUID {
			return ""
		} else {
			return utils.GenerateId()
		}
	}

	var getRuleParent func(ruleName string) (*api.RuleParent, error)
	getRuleName := func(rule *ast.Rule) string {
		return rule.Head.Name.String()
	}
	// always returns RuleChild even if "rule.Else" is present
	getRuleChild := func(rule *ast.Rule) (*api.RuleChild, error) {
		var statements []api.RuleStatement
		for i, expr := range rule.Body {
			var nextExpr *ast.Expr
			if i+1 < len(rule.Body) {
				nextExpr = rule.Body[i+1]
			}
			if expr.Generated {
				continue
			}
			var irStmts []ir.Stmt
			firstFound := false
			for _, block := range ruleNameToFunc[rule.Head.Name.String()].Blocks {
				flatIrStmts := flattenIrStmts(block.Stmts)
				for _, stmt := range flatIrStmts {
					if !firstFound {
						if stmt.GetLocation().Row >= expr.Location.Row && stmt.GetLocation().Col >= expr.Location.Col {
							firstFound = true
						}
					} else {
						if nextExpr != nil && (stmt.GetLocation().Row >= nextExpr.Location.Row && stmt.GetLocation().Col >= nextExpr.Location.Col) {
							break
						}
					}
					if firstFound {
						irStmts = append(irStmts, stmt)
					}
				}
				if firstFound {
					break
				}
			}
			if len(irStmts) == 0 {
				return nil, fmt.Errorf("len(irStmts) == 0")
			}

			var deps []api.RuleStatementDependenciesItem
			for _, irStmt := range irStmts {
				switch irStmt := irStmt.(type) {
				case *ir.CallStmt:
					ruleName := funcNameToRuleName[irStmt.Func]
					if ruleName == "" {
						// built-in function or something
						continue
					}
					ruleParent, err := getRuleParent(ruleName)
					if err != nil {
						return nil, err
					}
					deps = append(deps, api.RuleStatementDependenciesItem{
						Type:       api.RuleParentRuleStatementDependenciesItem,
						RuleParent: *ruleParent,
					})
				// todo: support like `data.foo.bar`, `import data.foo.bar`
				case *ir.DotStmt:
					switch src := irStmt.Source.Value.(type) {
					case *ir.Local:
						baseDocument := ""
						if *src == ir.Input {
							baseDocument = "input"
						} else if *src == ir.Data {
							baseDocument = "data"
						}
						if baseDocument != "" {
							switch key := irStmt.Key.Value.(type) {
							case *ir.StringIndex:
								keyString := moduleIr.Static.Strings[*key].Value
								deps = append(deps, api.RuleStatementDependenciesItem{
									Type:   api.StringRuleStatementDependenciesItem,
									String: fmt.Sprintf("%s.%s", baseDocument, keyString),
								})
							}
						}
					}
				}
			}

			statements = append(statements, api.RuleStatement{
				Name:         string(expr.Location.Text),
				UID:          maybeGenerateUID(),
				Dependencies: deps,
			})
		}
		value := ""
		if rule.Head.Value != nil {
			value = rule.Head.Value.String()
		}
		result := &api.RuleChild{
			Name:       fmt.Sprintf("%s:%d", getRuleName(rule), rule.Head.Location.Row),
			UID:        maybeGenerateUID(),
			Type:       api.RuleChildTypeChild,
			Value:      value,
			Statements: statements,
		}
		return result, nil
	}

	// must pass non-default rules
	getRuleParentChild := func(rule *ast.Rule) (*api.RuleParentChildrenItem, error) {
		var ruleChildElseChildren []api.RuleChild
		current := rule.Else
		for current != nil {
			ruleChild, err := getRuleChild(current)
			if err != nil {
				return nil, err
			}
			ruleChildElseChildren = append(ruleChildElseChildren, *ruleChild)
			current = current.Else
		}

		ruleChild, err := getRuleChild(rule)
		if err != nil {
			return nil, err
		}
		if ruleChildElseChildren != nil {
			ruleChildElseChildren = append(ruleChildElseChildren, *ruleChild)
			return &api.RuleParentChildrenItem{
				Type: api.RuleChildElseRuleParentChildrenItem,
				RuleChildElse: api.RuleChildElse{
					Name:     getRuleName(rule),
					UID:      maybeGenerateUID(),
					Type:     api.RuleChildElseTypeChildElse,
					Children: ruleChildElseChildren,
				},
			}, nil
		} else {
			return &api.RuleParentChildrenItem{
				Type:      api.RuleChildRuleParentChildrenItem,
				RuleChild: *ruleChild,
			}, nil
		}
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
	getRuleParent = func(ruleName string) (*api.RuleParent, error) {
		var nonDefaultRules []*ast.Rule
		var defaultRule *ast.Rule
		for _, rule := range moduleAst.Rules {
			if rule.Head.Name.String() == ruleName {
				if rule.Default {
					if defaultRule != nil {
						return nil, fmt.Errorf("multiple default rules found: %s", ruleName)
					}
					defaultRule = rule
				} else {
					nonDefaultRules = append(nonDefaultRules, rule)
				}
			}
		}
		if len(nonDefaultRules) == 0 {
			return nil, fmt.Errorf("rule not found: %s", ruleName)
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
			UID:      maybeGenerateUID(),
			Type:     api.RuleParentTypeParent,
			Default:  defaultRuleValue,
			Children: children,
		}
		return result, nil
	}

	result, err := getRuleParent(entrypoint)
	if err != nil {
		return nil, err
	}

	return result, nil
}

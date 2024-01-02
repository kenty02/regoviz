package main

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
	"strings"
)

func compileRego(moduleCode string) (*ast.Module, error) {
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

// from opa-explorer
func plan(ctx context.Context, rego string, print bool) (*ir.Policy, error) {
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

func planAsText(ctx context.Context, rego string, print bool) (string, error) {
	policy, err := plan(ctx, rego, print)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	if err := ir.Pretty(&buf, policy); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getDepTreeMap(policy *ir.Policy) map[string][]string {
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

func getDepTreePretty(policy *ir.Policy) string {
	depTreeMap := getDepTreeMap(policy)
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
func getMermaidFlowchart(policy *ir.Policy) string {
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

func getMermaidUrl(code string, edit bool) (string, error) {
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
	varLineNum int
	varName    string
}

type FixVarCommand struct {
	varLineNum int
	varName    string
	varValue   string
}

// varトークンにマッチする正規表現パターン
// \bは単語の境界を意味し、hogeが単独のトークンとして現れる場合にのみマッチする
var varPattern = `\b[A-Za-z_][A-Za-z_0-9]*\b`
var varRe = regexp.MustCompile(varPattern)

// 正規表現パターン: VTBEGINとVTENDの間にあるvarNameとvarValueを抽出
var vtPattern = `VTBEGIN ([A-Za-z_][A-Za-z_0-9]*) (.+) VTEND`
var vtRe = regexp.MustCompile(vtPattern)

func regoVarTrace(code, query string, input, data map[string]interface{}, commands []interface{}) (string, error) {
	lines := strings.Split(code, "\n")
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
		index := command.varLineNum - 1
		line := lines[index]

		// マッチがあるかどうか調べる
		matches := varRe.FindAllString(line, -1)
		// check if varName is in matches
		varNameFound := false
		for _, match := range matches {
			if match == command.varName {
				varNameFound = true
				// random suffixed var name
				suffix := rand.Intn(100000)
				newVarName := fmt.Sprintf("%s_fixed_%d", command.varName, suffix)
				// そのマッチを置換
				var err error
				line, err = replaceTokenInLine(line, command.varName, newVarName) // TODO: BNFと一貫性のある置換
				if err != nil {
					return "", err
				}
			}
		}
		if !varNameFound {
			return "", fmt.Errorf("var %s not found in line %d", command.varName, command.varLineNum)
		}
		line = line + fmt.Sprintf("; %s:=%s", command.varName, command.varValue)
		cis := []CodeInject{
			{command.varLineNum, line, true}}
		modifiedCode = injectCode(modifiedCode, cis)
		// just test if it works
		_, _, err := evalRegoWithPrint(modifiedCode, query, input, nil)
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprintf("変数%sの値を%sに固定しました。", command.varName, command.varValue))
	}
	for _, command := range showVarsCommands {
		// 行番号は1から始まるが、スライスのインデックスは0から始まるため調整
		index := command.varLineNum - 1
		line := lines[index]

		// マッチがあるかどうか調べる
		matches := varRe.FindAllString(line, -1)
		// check if varName is in matches
		varNameFound := false
		for _, match := range matches {
			if match == command.varName {
				varNameFound = true
				break
			}
		}
		if !varNameFound {
			return "", fmt.Errorf("var %s not found in line %d", command.varName, command.varLineNum)
		}

		cis := []CodeInject{
			{command.varLineNum, fmt.Sprintf(";print(\"VTBEGIN %s\", %s, \"VTEND\");false", command.varName, command.varName), false},
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
				fmt.Printf("varName: %s, varValue: %s\n", varName, varValue)

				if varName != command.varName {
					return "", fmt.Errorf("expected varName: %s, actual varName: %s", command.varName, varName)
				}
				varSet.Add(varValue)
			}
		}

		cis2 := []CodeInject{
			{command.varLineNum, fmt.Sprintf(";print(\"VTBEGIN %s\", %s, \"VTEND\")", command.varName, command.varName), false},
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
				fmt.Printf("varName: %s, varValue: %s\n", varName, varValue)

				if varName != command.varName {
					return "", fmt.Errorf("expected varName: %s, actual varName: %s", command.varName, varName)
				}
				actualVarSet.Add(varValue)
			}
		}

		if !varSet.IsEmpty() {
			sb.WriteString(fmt.Sprintf("変数%sは%sの値を取り得ます。", command.varName, varSet.ToSlice()))
			if !actualVarSet.IsEmpty() {
				sb.WriteString(fmt.Sprintf("評価パスで実際に代入された値の集合は%sです。", actualVarSet.ToSlice()))
			}
		} else {
			sb.WriteString(fmt.Sprintf("変数%sは値を取り得ませんでした。", command.varName))
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"
)

type AstGetOK struct {
	// The generated AST.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *AstGetOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *AstGetOK) SetResult(val string) {
	s.Result = val
}

type AstPrettyGetOK struct {
	// The generated AST.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *AstPrettyGetOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *AstPrettyGetOK) SetResult(val string) {
	s.Result = val
}

type BearerAuth struct {
	Token string
}

// GetToken returns the value of Token.
func (s *BearerAuth) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *BearerAuth) SetToken(val string) {
	s.Token = val
}

type CallTreeAvailableEntrypointsGetOK struct {
	Entrypoints []string `json:"entrypoints"`
}

// GetEntrypoints returns the value of Entrypoints.
func (s *CallTreeAvailableEntrypointsGetOK) GetEntrypoints() []string {
	return s.Entrypoints
}

// SetEntrypoints sets the value of Entrypoints.
func (s *CallTreeAvailableEntrypointsGetOK) SetEntrypoints(val []string) {
	s.Entrypoints = val
}

type CallTreeGetOK struct {
	Entrypoint RuleParent `json:"entrypoint"`
	Steps      []EvalStep `json:"steps"`
}

// GetEntrypoint returns the value of Entrypoint.
func (s *CallTreeGetOK) GetEntrypoint() RuleParent {
	return s.Entrypoint
}

// GetSteps returns the value of Steps.
func (s *CallTreeGetOK) GetSteps() []EvalStep {
	return s.Steps
}

// SetEntrypoint sets the value of Entrypoint.
func (s *CallTreeGetOK) SetEntrypoint(val RuleParent) {
	s.Entrypoint = val
}

// SetSteps sets the value of Steps.
func (s *CallTreeGetOK) SetSteps(val []EvalStep) {
	s.Steps = val
}

type DepTreeTextGetOK struct {
	// The generated dependency tree.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *DepTreeTextGetOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *DepTreeTextGetOK) SetResult(val string) {
	s.Result = val
}

// Ref: #/components/schemas/EvalStep
type EvalStep struct {
	Index         int    `json:"index"`
	Message       string `json:"message"`
	TargetNodeUid string `json:"targetNodeUid"`
}

// GetIndex returns the value of Index.
func (s *EvalStep) GetIndex() int {
	return s.Index
}

// GetMessage returns the value of Message.
func (s *EvalStep) GetMessage() string {
	return s.Message
}

// GetTargetNodeUid returns the value of TargetNodeUid.
func (s *EvalStep) GetTargetNodeUid() string {
	return s.TargetNodeUid
}

// SetIndex sets the value of Index.
func (s *EvalStep) SetIndex(val int) {
	s.Index = val
}

// SetMessage sets the value of Message.
func (s *EvalStep) SetMessage(val string) {
	s.Message = val
}

// SetTargetNodeUid sets the value of TargetNodeUid.
func (s *EvalStep) SetTargetNodeUid(val string) {
	s.TargetNodeUid = val
}

type FlowchartGetOK struct {
	// The generated flowchart mermaid url.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *FlowchartGetOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *FlowchartGetOK) SetResult(val string) {
	s.Result = val
}

type IrGetOK struct {
	// The generated IR.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *IrGetOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *IrGetOK) SetResult(val string) {
	s.Result = val
}

// Ref: #/components/schemas/NodeLocation
type NodeLocation struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// GetRow returns the value of Row.
func (s *NodeLocation) GetRow() int {
	return s.Row
}

// GetCol returns the value of Col.
func (s *NodeLocation) GetCol() int {
	return s.Col
}

// SetRow sets the value of Row.
func (s *NodeLocation) SetRow(val int) {
	s.Row = val
}

// SetCol sets the value of Col.
func (s *NodeLocation) SetCol(val int) {
	s.Col = val
}

// NewOptBool returns new OptBool with value set to v.
func NewOptBool(v bool) OptBool {
	return OptBool{
		Value: v,
		Set:   true,
	}
}

// OptBool is optional bool.
type OptBool struct {
	Value bool
	Set   bool
}

// IsSet returns true if OptBool was set.
func (o OptBool) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptBool) Reset() {
	var v bool
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptBool) SetTo(v bool) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptBool) Get() (v bool, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptBool) Or(d bool) bool {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptNodeLocation returns new OptNodeLocation with value set to v.
func NewOptNodeLocation(v NodeLocation) OptNodeLocation {
	return OptNodeLocation{
		Value: v,
		Set:   true,
	}
}

// OptNodeLocation is optional NodeLocation.
type OptNodeLocation struct {
	Value NodeLocation
	Set   bool
}

// IsSet returns true if OptNodeLocation was set.
func (o OptNodeLocation) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNodeLocation) Reset() {
	var v NodeLocation
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptNodeLocation) SetTo(v NodeLocation) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNodeLocation) Get() (v NodeLocation, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNodeLocation) Or(d NodeLocation) NodeLocation {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Merged schema.
// Ref: #/components/schemas/RuleChild
type RuleChild struct {
	Name       string          `json:"name"`
	UID        string          `json:"uid"`
	Location   OptNodeLocation `json:"location"`
	Type       RuleChildType   `json:"type"`
	Value      string          `json:"value"`
	Statements []RuleStatement `json:"statements"`
}

// GetName returns the value of Name.
func (s *RuleChild) GetName() string {
	return s.Name
}

// GetUID returns the value of UID.
func (s *RuleChild) GetUID() string {
	return s.UID
}

// GetLocation returns the value of Location.
func (s *RuleChild) GetLocation() OptNodeLocation {
	return s.Location
}

// GetType returns the value of Type.
func (s *RuleChild) GetType() RuleChildType {
	return s.Type
}

// GetValue returns the value of Value.
func (s *RuleChild) GetValue() string {
	return s.Value
}

// GetStatements returns the value of Statements.
func (s *RuleChild) GetStatements() []RuleStatement {
	return s.Statements
}

// SetName sets the value of Name.
func (s *RuleChild) SetName(val string) {
	s.Name = val
}

// SetUID sets the value of UID.
func (s *RuleChild) SetUID(val string) {
	s.UID = val
}

// SetLocation sets the value of Location.
func (s *RuleChild) SetLocation(val OptNodeLocation) {
	s.Location = val
}

// SetType sets the value of Type.
func (s *RuleChild) SetType(val RuleChildType) {
	s.Type = val
}

// SetValue sets the value of Value.
func (s *RuleChild) SetValue(val string) {
	s.Value = val
}

// SetStatements sets the value of Statements.
func (s *RuleChild) SetStatements(val []RuleStatement) {
	s.Statements = val
}

// Merged schema.
// Ref: #/components/schemas/RuleChildElse
type RuleChildElse struct {
	Name     string            `json:"name"`
	UID      string            `json:"uid"`
	Location OptNodeLocation   `json:"location"`
	Type     RuleChildElseType `json:"type"`
	Children []RuleChild       `json:"children"`
}

// GetName returns the value of Name.
func (s *RuleChildElse) GetName() string {
	return s.Name
}

// GetUID returns the value of UID.
func (s *RuleChildElse) GetUID() string {
	return s.UID
}

// GetLocation returns the value of Location.
func (s *RuleChildElse) GetLocation() OptNodeLocation {
	return s.Location
}

// GetType returns the value of Type.
func (s *RuleChildElse) GetType() RuleChildElseType {
	return s.Type
}

// GetChildren returns the value of Children.
func (s *RuleChildElse) GetChildren() []RuleChild {
	return s.Children
}

// SetName sets the value of Name.
func (s *RuleChildElse) SetName(val string) {
	s.Name = val
}

// SetUID sets the value of UID.
func (s *RuleChildElse) SetUID(val string) {
	s.UID = val
}

// SetLocation sets the value of Location.
func (s *RuleChildElse) SetLocation(val OptNodeLocation) {
	s.Location = val
}

// SetType sets the value of Type.
func (s *RuleChildElse) SetType(val RuleChildElseType) {
	s.Type = val
}

// SetChildren sets the value of Children.
func (s *RuleChildElse) SetChildren(val []RuleChild) {
	s.Children = val
}

type RuleChildElseType string

const (
	RuleChildElseTypeChildElse RuleChildElseType = "child-else"
)

// AllValues returns all RuleChildElseType values.
func (RuleChildElseType) AllValues() []RuleChildElseType {
	return []RuleChildElseType{
		RuleChildElseTypeChildElse,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RuleChildElseType) MarshalText() ([]byte, error) {
	switch s {
	case RuleChildElseTypeChildElse:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RuleChildElseType) UnmarshalText(data []byte) error {
	switch RuleChildElseType(data) {
	case RuleChildElseTypeChildElse:
		*s = RuleChildElseTypeChildElse
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

type RuleChildType string

const (
	RuleChildTypeChild RuleChildType = "child"
)

// AllValues returns all RuleChildType values.
func (RuleChildType) AllValues() []RuleChildType {
	return []RuleChildType{
		RuleChildTypeChild,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RuleChildType) MarshalText() ([]byte, error) {
	switch s {
	case RuleChildTypeChild:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RuleChildType) UnmarshalText(data []byte) error {
	switch RuleChildType(data) {
	case RuleChildTypeChild:
		*s = RuleChildTypeChild
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Merged schema.
// Ref: #/components/schemas/RuleParent
type RuleParent struct {
	Name     string                   `json:"name"`
	UID      string                   `json:"uid"`
	Location OptNodeLocation          `json:"location"`
	Type     RuleParentType           `json:"type"`
	Default  string                   `json:"default"`
	Children []RuleParentChildrenItem `json:"children"`
	Ref      string                   `json:"ref"`
}

// GetName returns the value of Name.
func (s *RuleParent) GetName() string {
	return s.Name
}

// GetUID returns the value of UID.
func (s *RuleParent) GetUID() string {
	return s.UID
}

// GetLocation returns the value of Location.
func (s *RuleParent) GetLocation() OptNodeLocation {
	return s.Location
}

// GetType returns the value of Type.
func (s *RuleParent) GetType() RuleParentType {
	return s.Type
}

// GetDefault returns the value of Default.
func (s *RuleParent) GetDefault() string {
	return s.Default
}

// GetChildren returns the value of Children.
func (s *RuleParent) GetChildren() []RuleParentChildrenItem {
	return s.Children
}

// GetRef returns the value of Ref.
func (s *RuleParent) GetRef() string {
	return s.Ref
}

// SetName sets the value of Name.
func (s *RuleParent) SetName(val string) {
	s.Name = val
}

// SetUID sets the value of UID.
func (s *RuleParent) SetUID(val string) {
	s.UID = val
}

// SetLocation sets the value of Location.
func (s *RuleParent) SetLocation(val OptNodeLocation) {
	s.Location = val
}

// SetType sets the value of Type.
func (s *RuleParent) SetType(val RuleParentType) {
	s.Type = val
}

// SetDefault sets the value of Default.
func (s *RuleParent) SetDefault(val string) {
	s.Default = val
}

// SetChildren sets the value of Children.
func (s *RuleParent) SetChildren(val []RuleParentChildrenItem) {
	s.Children = val
}

// SetRef sets the value of Ref.
func (s *RuleParent) SetRef(val string) {
	s.Ref = val
}

// RuleParentChildrenItem represents sum type.
type RuleParentChildrenItem struct {
	Type          RuleParentChildrenItemType // switch on this field
	RuleChild     RuleChild
	RuleChildElse RuleChildElse
}

// RuleParentChildrenItemType is oneOf type of RuleParentChildrenItem.
type RuleParentChildrenItemType string

// Possible values for RuleParentChildrenItemType.
const (
	RuleChildRuleParentChildrenItem     RuleParentChildrenItemType = "RuleChild"
	RuleChildElseRuleParentChildrenItem RuleParentChildrenItemType = "RuleChildElse"
)

// IsRuleChild reports whether RuleParentChildrenItem is RuleChild.
func (s RuleParentChildrenItem) IsRuleChild() bool { return s.Type == RuleChildRuleParentChildrenItem }

// IsRuleChildElse reports whether RuleParentChildrenItem is RuleChildElse.
func (s RuleParentChildrenItem) IsRuleChildElse() bool {
	return s.Type == RuleChildElseRuleParentChildrenItem
}

// SetRuleChild sets RuleParentChildrenItem to RuleChild.
func (s *RuleParentChildrenItem) SetRuleChild(v RuleChild) {
	s.Type = RuleChildRuleParentChildrenItem
	s.RuleChild = v
}

// GetRuleChild returns RuleChild and true boolean if RuleParentChildrenItem is RuleChild.
func (s RuleParentChildrenItem) GetRuleChild() (v RuleChild, ok bool) {
	if !s.IsRuleChild() {
		return v, false
	}
	return s.RuleChild, true
}

// NewRuleChildRuleParentChildrenItem returns new RuleParentChildrenItem from RuleChild.
func NewRuleChildRuleParentChildrenItem(v RuleChild) RuleParentChildrenItem {
	var s RuleParentChildrenItem
	s.SetRuleChild(v)
	return s
}

// SetRuleChildElse sets RuleParentChildrenItem to RuleChildElse.
func (s *RuleParentChildrenItem) SetRuleChildElse(v RuleChildElse) {
	s.Type = RuleChildElseRuleParentChildrenItem
	s.RuleChildElse = v
}

// GetRuleChildElse returns RuleChildElse and true boolean if RuleParentChildrenItem is RuleChildElse.
func (s RuleParentChildrenItem) GetRuleChildElse() (v RuleChildElse, ok bool) {
	if !s.IsRuleChildElse() {
		return v, false
	}
	return s.RuleChildElse, true
}

// NewRuleChildElseRuleParentChildrenItem returns new RuleParentChildrenItem from RuleChildElse.
func NewRuleChildElseRuleParentChildrenItem(v RuleChildElse) RuleParentChildrenItem {
	var s RuleParentChildrenItem
	s.SetRuleChildElse(v)
	return s
}

type RuleParentType string

const (
	RuleParentTypeParent RuleParentType = "parent"
)

// AllValues returns all RuleParentType values.
func (RuleParentType) AllValues() []RuleParentType {
	return []RuleParentType{
		RuleParentTypeParent,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RuleParentType) MarshalText() ([]byte, error) {
	switch s {
	case RuleParentTypeParent:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RuleParentType) UnmarshalText(data []byte) error {
	switch RuleParentType(data) {
	case RuleParentTypeParent:
		*s = RuleParentTypeParent
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Merged schema.
// Ref: #/components/schemas/RuleStatement
type RuleStatement struct {
	Name         string                          `json:"name"`
	UID          string                          `json:"uid"`
	Location     OptNodeLocation                 `json:"location"`
	Dependencies []RuleStatementDependenciesItem `json:"dependencies"`
}

// GetName returns the value of Name.
func (s *RuleStatement) GetName() string {
	return s.Name
}

// GetUID returns the value of UID.
func (s *RuleStatement) GetUID() string {
	return s.UID
}

// GetLocation returns the value of Location.
func (s *RuleStatement) GetLocation() OptNodeLocation {
	return s.Location
}

// GetDependencies returns the value of Dependencies.
func (s *RuleStatement) GetDependencies() []RuleStatementDependenciesItem {
	return s.Dependencies
}

// SetName sets the value of Name.
func (s *RuleStatement) SetName(val string) {
	s.Name = val
}

// SetUID sets the value of UID.
func (s *RuleStatement) SetUID(val string) {
	s.UID = val
}

// SetLocation sets the value of Location.
func (s *RuleStatement) SetLocation(val OptNodeLocation) {
	s.Location = val
}

// SetDependencies sets the value of Dependencies.
func (s *RuleStatement) SetDependencies(val []RuleStatementDependenciesItem) {
	s.Dependencies = val
}

// RuleStatementDependenciesItem represents sum type.
type RuleStatementDependenciesItem struct {
	Type       RuleStatementDependenciesItemType // switch on this field
	RuleParent RuleParent
	String     string
}

// RuleStatementDependenciesItemType is oneOf type of RuleStatementDependenciesItem.
type RuleStatementDependenciesItemType string

// Possible values for RuleStatementDependenciesItemType.
const (
	RuleParentRuleStatementDependenciesItem RuleStatementDependenciesItemType = "RuleParent"
	StringRuleStatementDependenciesItem     RuleStatementDependenciesItemType = "string"
)

// IsRuleParent reports whether RuleStatementDependenciesItem is RuleParent.
func (s RuleStatementDependenciesItem) IsRuleParent() bool {
	return s.Type == RuleParentRuleStatementDependenciesItem
}

// IsString reports whether RuleStatementDependenciesItem is string.
func (s RuleStatementDependenciesItem) IsString() bool {
	return s.Type == StringRuleStatementDependenciesItem
}

// SetRuleParent sets RuleStatementDependenciesItem to RuleParent.
func (s *RuleStatementDependenciesItem) SetRuleParent(v RuleParent) {
	s.Type = RuleParentRuleStatementDependenciesItem
	s.RuleParent = v
}

// GetRuleParent returns RuleParent and true boolean if RuleStatementDependenciesItem is RuleParent.
func (s RuleStatementDependenciesItem) GetRuleParent() (v RuleParent, ok bool) {
	if !s.IsRuleParent() {
		return v, false
	}
	return s.RuleParent, true
}

// NewRuleParentRuleStatementDependenciesItem returns new RuleStatementDependenciesItem from RuleParent.
func NewRuleParentRuleStatementDependenciesItem(v RuleParent) RuleStatementDependenciesItem {
	var s RuleStatementDependenciesItem
	s.SetRuleParent(v)
	return s
}

// SetString sets RuleStatementDependenciesItem to string.
func (s *RuleStatementDependenciesItem) SetString(v string) {
	s.Type = StringRuleStatementDependenciesItem
	s.String = v
}

// GetString returns string and true boolean if RuleStatementDependenciesItem is string.
func (s RuleStatementDependenciesItem) GetString() (v string, ok bool) {
	if !s.IsString() {
		return v, false
	}
	return s.String, true
}

// NewStringRuleStatementDependenciesItem returns new RuleStatementDependenciesItem from string.
func NewStringRuleStatementDependenciesItem(v string) RuleStatementDependenciesItem {
	var s RuleStatementDependenciesItem
	s.SetString(v)
	return s
}

// OPA policy that can be used with this API.
// Ref: #/components/schemas/Sample
type Sample struct {
	// The name of the sample file.
	FileName string `json:"file_name"`
	// The content of the sample file.
	Content string `json:"content"`
	// List of input examples for the sample.
	InputExamples SampleInputExamples `json:"input_examples"`
	// List of data examples for the sample.
	DataExamples SampleDataExamples `json:"data_examples"`
	// List of query examples for the sample.
	QueryExamples SampleQueryExamples `json:"query_examples"`
}

// GetFileName returns the value of FileName.
func (s *Sample) GetFileName() string {
	return s.FileName
}

// GetContent returns the value of Content.
func (s *Sample) GetContent() string {
	return s.Content
}

// GetInputExamples returns the value of InputExamples.
func (s *Sample) GetInputExamples() SampleInputExamples {
	return s.InputExamples
}

// GetDataExamples returns the value of DataExamples.
func (s *Sample) GetDataExamples() SampleDataExamples {
	return s.DataExamples
}

// GetQueryExamples returns the value of QueryExamples.
func (s *Sample) GetQueryExamples() SampleQueryExamples {
	return s.QueryExamples
}

// SetFileName sets the value of FileName.
func (s *Sample) SetFileName(val string) {
	s.FileName = val
}

// SetContent sets the value of Content.
func (s *Sample) SetContent(val string) {
	s.Content = val
}

// SetInputExamples sets the value of InputExamples.
func (s *Sample) SetInputExamples(val SampleInputExamples) {
	s.InputExamples = val
}

// SetDataExamples sets the value of DataExamples.
func (s *Sample) SetDataExamples(val SampleDataExamples) {
	s.DataExamples = val
}

// SetQueryExamples sets the value of QueryExamples.
func (s *Sample) SetQueryExamples(val SampleQueryExamples) {
	s.QueryExamples = val
}

// List of data examples for the sample.
type SampleDataExamples struct {
	// The default data for the sample. Can be empty.
	Default         string `json:"default"`
	AdditionalProps SampleDataExamplesAdditional
}

// GetDefault returns the value of Default.
func (s *SampleDataExamples) GetDefault() string {
	return s.Default
}

// GetAdditionalProps returns the value of AdditionalProps.
func (s *SampleDataExamples) GetAdditionalProps() SampleDataExamplesAdditional {
	return s.AdditionalProps
}

// SetDefault sets the value of Default.
func (s *SampleDataExamples) SetDefault(val string) {
	s.Default = val
}

// SetAdditionalProps sets the value of AdditionalProps.
func (s *SampleDataExamples) SetAdditionalProps(val SampleDataExamplesAdditional) {
	s.AdditionalProps = val
}

type SampleDataExamplesAdditional map[string]string

func (s *SampleDataExamplesAdditional) init() SampleDataExamplesAdditional {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

// List of input examples for the sample.
type SampleInputExamples struct {
	// The default input for the sample. Can be empty.
	Default         string `json:"default"`
	AdditionalProps SampleInputExamplesAdditional
}

// GetDefault returns the value of Default.
func (s *SampleInputExamples) GetDefault() string {
	return s.Default
}

// GetAdditionalProps returns the value of AdditionalProps.
func (s *SampleInputExamples) GetAdditionalProps() SampleInputExamplesAdditional {
	return s.AdditionalProps
}

// SetDefault sets the value of Default.
func (s *SampleInputExamples) SetDefault(val string) {
	s.Default = val
}

// SetAdditionalProps sets the value of AdditionalProps.
func (s *SampleInputExamples) SetAdditionalProps(val SampleInputExamplesAdditional) {
	s.AdditionalProps = val
}

type SampleInputExamplesAdditional map[string]string

func (s *SampleInputExamplesAdditional) init() SampleInputExamplesAdditional {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

// List of query examples for the sample.
type SampleQueryExamples struct {
	// The default query for the sample. Can be empty.
	Default         string `json:"default"`
	AdditionalProps SampleQueryExamplesAdditional
}

// GetDefault returns the value of Default.
func (s *SampleQueryExamples) GetDefault() string {
	return s.Default
}

// GetAdditionalProps returns the value of AdditionalProps.
func (s *SampleQueryExamples) GetAdditionalProps() SampleQueryExamplesAdditional {
	return s.AdditionalProps
}

// SetDefault sets the value of Default.
func (s *SampleQueryExamples) SetDefault(val string) {
	s.Default = val
}

// SetAdditionalProps sets the value of AdditionalProps.
func (s *SampleQueryExamples) SetAdditionalProps(val SampleQueryExamplesAdditional) {
	s.AdditionalProps = val
}

type SampleQueryExamplesAdditional map[string]string

func (s *SampleQueryExamplesAdditional) init() SampleQueryExamplesAdditional {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

type VarTracePostOK struct {
	// The output of variable trace.
	Result string `json:"result"`
}

// GetResult returns the value of Result.
func (s *VarTracePostOK) GetResult() string {
	return s.Result
}

// SetResult sets the value of Result.
func (s *VarTracePostOK) SetResult(val string) {
	s.Result = val
}

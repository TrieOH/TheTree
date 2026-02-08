package goauth

type Condition struct {
	And *[]Condition `json:"and,omitempty"`
	Or  *[]Condition `json:"or,omitempty"`
	Not *Condition   `json:"not,omitempty"`

	// For leaf conditions (only one of these should be set)
	Path       string      `json:"path,omitempty"`
	Op         string      `json:"op,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Ref        string      `json:"ref,omitempty"`
	Margin     string      `json:"margin,omitempty"`
	Field      string      `json:"field,omitempty"`
	FieldStart string      `json:"field_start,omitempty"`
	FieldEnd   string      `json:"field_end,omitempty"`

	// Permission check specific
	Action  string  `json:"action,omitempty"`
	Object  string  `json:"object,omitempty"`
	ScopeID *string `json:"scope_id,omitempty"`
}

// ConditionBuilder is the definitive end of a condition definition.
type ConditionBuilder interface {
	Build() Condition
}

// ConditionFactory is the entry point for building complex conditions.
// It provides a structured way to nest conditions with proper indentation.
type ConditionFactory struct{}

// NewCondition returns a new ConditionFactory to start building conditions.
func NewCondition() ConditionFactory {
	return ConditionFactory{}
}

func (f ConditionFactory) And(conds ...ConditionBuilder) ConditionBuilder {
	c := make([]Condition, len(conds))
	for i, cb := range conds {
		c[i] = cb.Build()
	}
	return &baseConditionBuilder{cond: Condition{And: &c}}
}

func (f ConditionFactory) Or(conds ...ConditionBuilder) ConditionBuilder {
	c := make([]Condition, len(conds))
	for i, cb := range conds {
		c[i] = cb.Build()
	}
	return &baseConditionBuilder{cond: Condition{Or: &c}}
}

func (f ConditionFactory) Not(cond ConditionBuilder) ConditionBuilder {
	c := cond.Build()
	return &baseConditionBuilder{cond: Condition{Not: &c}}
}

func (f ConditionFactory) Path(path string) PathPredicateBuilder {
	return &pathBuilder{path: path}
}

func (f ConditionFactory) Field(field string) TemporalGraceBuilder {
	return &fieldBuilder{field: field}
}

func (f ConditionFactory) Fields(start, end string) GraceDurationBuilder {
	return &graceDurationBuilder{start: start, end: end}
}

type PathPredicateBuilder interface {
	Eq(val interface{}) ConditionBuilder
	Neq(val interface{}) ConditionBuilder
	Gt(val interface{}) ConditionBuilder
	Gte(val interface{}) ConditionBuilder
	Lt(val interface{}) ConditionBuilder
	Lte(val interface{}) ConditionBuilder
	StartsWith(val string) ConditionBuilder
	EndsWith(val string) ConditionBuilder
	Contains(val interface{}) ConditionBuilder
	ContainsAll(val interface{}) ConditionBuilder
	ContainsAny(val interface{}) ConditionBuilder
	Matches(regex string) ConditionBuilder
	In(val interface{}) ConditionBuilder
	Exists() ConditionBuilder
	RefEq(ref string) ConditionBuilder
	RefNeq(ref string) ConditionBuilder
}

type TemporalGraceBuilder interface {
	GraceBefore(margin string) ConditionBuilder
	GraceAfter(margin string) ConditionBuilder
	GraceAround(margin string) ConditionBuilder
}

type GraceDurationBuilder interface {
	GraceDuration(margin string) ConditionBuilder
}

type baseConditionBuilder struct {
	cond Condition
}

func (b *baseConditionBuilder) Build() Condition {
	return b.cond
}

type pathBuilder struct {
	path string
}

func (p *pathBuilder) build(op string, val interface{}, ref string) ConditionBuilder {
	return &baseConditionBuilder{cond: Condition{Path: p.path, Op: op, Value: val, Ref: ref}}
}

func (p *pathBuilder) Eq(val interface{}) ConditionBuilder          { return p.build("eq", val, "") }
func (p *pathBuilder) Neq(val interface{}) ConditionBuilder         { return p.build("neq", val, "") }
func (p *pathBuilder) Gt(val interface{}) ConditionBuilder          { return p.build("gt", val, "") }
func (p *pathBuilder) Gte(val interface{}) ConditionBuilder         { return p.build("gte", val, "") }
func (p *pathBuilder) Lt(val interface{}) ConditionBuilder          { return p.build("lt", val, "") }
func (p *pathBuilder) Lte(val interface{}) ConditionBuilder         { return p.build("lte", val, "") }
func (p *pathBuilder) StartsWith(val string) ConditionBuilder       { return p.build("startsWith", val, "") }
func (p *pathBuilder) EndsWith(val string) ConditionBuilder         { return p.build("endsWith", val, "") }
func (p *pathBuilder) Contains(val interface{}) ConditionBuilder    { return p.build("contains", val, "") }
func (p *pathBuilder) ContainsAll(val interface{}) ConditionBuilder { return p.build("containsAll", val, "") }
func (p *pathBuilder) ContainsAny(val interface{}) ConditionBuilder { return p.build("containsAny", val, "") }
func (p *pathBuilder) Matches(regex string) ConditionBuilder        { return p.build("matches", regex, "") }
func (p *pathBuilder) In(val interface{}) ConditionBuilder          { return p.build("in", val, "") }
func (p *pathBuilder) Exists() ConditionBuilder                    { return p.build("exists", nil, "") }
func (p *pathBuilder) RefEq(ref string) ConditionBuilder            { return p.build("eq", nil, ref) }
func (p *pathBuilder) RefNeq(ref string) ConditionBuilder           { return p.build("neq", nil, ref) }

type fieldBuilder struct {
	field string
}

func (f *fieldBuilder) build(op, margin string) ConditionBuilder {
	return &baseConditionBuilder{cond: Condition{Field: f.field, Op: op, Margin: margin}}
}

func (f *fieldBuilder) GraceBefore(margin string) ConditionBuilder { return f.build("grace_before", margin) }
func (f *fieldBuilder) GraceAfter(margin string) ConditionBuilder  { return f.build("grace_after", margin) }
func (f *fieldBuilder) GraceAround(margin string) ConditionBuilder { return f.build("grace_around", margin) }

type graceDurationBuilder struct {
	start string
	end   string
}

func (g *graceDurationBuilder) GraceDuration(margin string) ConditionBuilder {
	return &baseConditionBuilder{cond: Condition{FieldStart: g.start, FieldEnd: g.end, Op: "grace_duration", Margin: margin}}
}
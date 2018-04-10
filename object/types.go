package object

import (
	"strconv"
	"strings"
)

const (
	TYPE_NULL   = "NULL"
	TYPE_NUMBER = "NUMBER"
	TYPE_STRING = "STRING"
	TYPE_BOOL   = "BOOLEAN"
	TYPE_FILTER = "FILTER"
	TYPE_ARRAY  = "ARRAY"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{value: true}
	FALSE = &Boolean{value: false}
)

type Null struct{}

func (n *Null) Type() ObjectType   { return TYPE_NULL }
func (n *Null) Value() interface{} { return nil }
func (n *Null) Inspect() string    { return "" }

type Number struct {
	value float64
}

func (n *Number) Type() ObjectType   { return TYPE_NUMBER }
func (n *Number) Value() interface{} { return n.value }
func (n *Number) Inspect() string    { return strconv.FormatFloat(n.value, 'f', -1, 64) }

type String struct {
	value string
}

func (s *String) Type() ObjectType   { return TYPE_STRING }
func (s *String) Value() interface{} { return s.value }
func (s *String) Inspect() string    { return s.value }

type Boolean struct {
	value bool
}

func (b *Boolean) Type() ObjectType   { return TYPE_BOOL }
func (b *Boolean) Value() interface{} { return b.value }
func (b *Boolean) Inspect() string {
	if b.value {
		return "true"
	} else {
		return "false"
	}
}

type Filter struct {
	Name       string
	Parameters map[string]Object
}

func (f *Filter) Type() ObjectType   { return TYPE_FILTER }
func (f *Filter) Value() interface{} { return f.Name }
func (f *Filter) Inspect() string    { return f.Name }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType   { return TYPE_ARRAY }
func (a *Array) Value() interface{} { return nil } // TODO What to return here?
func (a *Array) Inspect() string {
	output := strings.Builder{}

	var parts []string

	for _, expr := range a.Elements {
		parts = append(parts, expr.Inspect())
	}

	output.WriteString("[")
	output.WriteString(strings.Join(parts, ","))
	output.WriteString("]")

	return output.String()
}

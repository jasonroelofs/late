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
	TYPE_HASH   = "HASH"
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

func (a *Array) Append(obj Object) {
	a.Elements = append(a.Elements, obj)
}

func (a *Array) Get(index int) Object {
	return a.Elements[index]
}

func (a *Array) Len() int {
	return len(a.Elements)
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

type Hash struct {
	// Elements is keyed off an interface{} type to work around
	// Go's strict map-key-semantics. We can't make
	// object.String("one") == object.String("two") so to make our
	// hash work, we drop down to the underlyng Value() as the key.
	// This also means you shouldn't try to make a Hash or Array a key of another
	// Hash, but why would you want to do such a thing anyway?
	elements map[interface{}]Object
}

func NewHash() *Hash {
	return &Hash{
		elements: make(map[interface{}]Object),
	}
}

func (h *Hash) Get(key Object) Object {
	value, ok := h.elements[key.Value()]
	if !ok {
		return NULL
	}

	return value
}

func (h *Hash) Set(key Object, value Object) {
	h.elements[key.Value()] = value
}

func (h *Hash) Type() ObjectType   { return TYPE_HASH }
func (h *Hash) Value() interface{} { return nil } // TODO?
func (h *Hash) Inspect() string {
	return "TODO"
}

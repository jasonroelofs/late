package object

import (
	"strconv"
)

const (
	OBJ_NULL   = "NULL"
	OBJ_NUMBER = "NUMBER"
	OBJ_STRING = "STRING"
	OBJ_BOOL   = "BOOLEAN"
	OBJ_FILTER = "FILTER"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{value: true}
	FALSE = &Boolean{value: false}
)

type Null struct{}

func (n *Null) Type() ObjectType   { return OBJ_NULL }
func (n *Null) Value() interface{} { return nil }
func (n *Null) Inspect() string    { return "null" }

type Number struct {
	value float64
}

func (n *Number) Type() ObjectType   { return OBJ_NUMBER }
func (n *Number) Value() interface{} { return n.value }
func (n *Number) Inspect() string    { return strconv.FormatFloat(n.value, 'f', -1, 64) }

type String struct {
	value string
}

func (s *String) Type() ObjectType   { return OBJ_STRING }
func (s *String) Value() interface{} { return s.value }
func (s *String) Inspect() string    { return s.value }

type Boolean struct {
	value bool
}

func (b *Boolean) Type() ObjectType   { return OBJ_BOOL }
func (b *Boolean) Value() interface{} { return b.value }
func (b *Boolean) Inspect() string {
	if b.value {
		return "true"
	} else {
		return "false"
	}
}

type Filter struct {
	Name string
}

func (f *Filter) Type() ObjectType   { return OBJ_FILTER }
func (f *Filter) Value() interface{} { return f.Name }
func (f *Filter) Inspect() string    { return f.Name }

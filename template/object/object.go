package object

import (
	"strconv"
)

type ObjectType string

const (
	OBJ_NULL   = "NULL"
	OBJ_NUMBER = "NUMBER"
	OBJ_STRING = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType { return OBJ_NULL }
func (n *Null) Inspect() string  { return "null" }

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return OBJ_NUMBER }
func (n *Number) Inspect() string  { return strconv.FormatFloat(n.Value, 'f', -1, 64) }

type String struct {
	Value string
}

func (i *String) Type() ObjectType { return OBJ_STRING }
func (i *String) Inspect() string  { return i.Value }

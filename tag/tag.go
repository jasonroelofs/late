package tag

import (
	"github.com/jasonroelofs/late/object"
)

// Tags are the work-horse of Liquid templating.
// Where values take data and put it in the final output, tags are where
// the restricted programming logic lives.
// Tags exist inside of {% %} designators and can be stand-alone or in blocks
// that close with a matching {% end %}

type Tag interface {
	Parse() *ParseConfig
	Eval(Environment, *ParseResult) object.Object
}

type Environment interface {
	Eval(Statement) object.Object
	Get(string) object.Object
	Set(string, interface{})
}

type Statement interface {
	String() string
}

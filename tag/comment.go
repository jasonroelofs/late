package tag

import (
	"github.com/jasonroelofs/late/object"
)

/**
 * The comment block tag will throw away all content inside of the block.
 *
 *   {% comment %}
 *   ...
 *   {% end %}
 *
 */
type Comment struct {
}

func (a *Comment) Block() bool {
	return true
}

func (a *Comment) Parse() []ParseRule {
	return []ParseRule{}
}

func (a *Comment) Eval(_ Environment, _ []object.Object, _ []Statement) object.Object {
	// Do nothing, all content is gone
	return object.NULL
}

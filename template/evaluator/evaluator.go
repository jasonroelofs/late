package evaluator

import (
	"strings"

	"github.com/jasonroelofs/late/template/ast"
)

func Eval(template *ast.Template) string {
	buffer := strings.Builder{}

	for _, node := range template.Nodes {
		buffer.WriteString(node.String())
	}

	return buffer.String()
}

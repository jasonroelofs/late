package evaluator

import (
	"strings"

	"github.com/jasonroelofs/late/template/ast"
)

func Eval(template *ast.Template) string {
	buffer := strings.Builder{}

	for _, stmt := range template.Statements {
		buffer.WriteString(stmt.String())
	}

	return buffer.String()
}

package evaluator

import (
	"strings"

	"github.com/jasonroelofs/late/template/ast"
)

func Eval(template *ast.Template) string {
	buffer := strings.Builder{}

	for _, statement := range template.Statements {
		buffer.WriteString(statement.TokenLiteral())
	}

	return buffer.String()
}

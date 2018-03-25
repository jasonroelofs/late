package evaluator

import (
	"github.com/jasonroelofs/late"
	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/tag"
	"github.com/jasonroelofs/late/template/ast"
)

type Evaluator struct {
	context  *context.Context
	template *ast.Template
}

func New(template *ast.Template, context *context.Context) *Evaluator {
	return &Evaluator{
		template: template,
		context:  context,
	}
}

func (e *Evaluator) Run() []object.Object {
	var objects []object.Object

	for _, statement := range e.template.Statements {
		result := e.Eval(statement)
		objects = append(objects, result)
	}

	return objects
}

func (e *Evaluator) Set(variable string, value interface{}) {
	e.context.Set(variable, value)
}

func (e *Evaluator) Get(variable string) object.Object {
	return e.context.Get(variable)
}

func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Top-level Statements
	case *ast.RawStatement:
		return object.New(node.String())

	case *ast.VariableStatement:
		return e.Eval(node.Expression)

	case *ast.TagStatement:
		return e.evalTagStatement(node)

	// Expressions
	case *ast.InfixExpression:
		left := e.Eval(node.Left)
		right := e.Eval(node.Right)
		return e.evalInfix(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := e.Eval(node.Right)
		return e.evalPrefix(node.Operator, right)

	case *ast.FilterExpression:
		input := e.Eval(node.Input)
		filter := e.Eval(node.Filter)
		return e.evalFilter(input, filter)

	// Literals
	case *ast.NumberLiteral:
		return object.New(node.Value)

	case *ast.BooleanLiteral:
		return convertBoolean(node.Value)

	case *ast.StringLiteral:
		return object.New(node.Value)

	case *ast.Identifier:
		return e.evalIdentifier(node.Value)

	case *ast.FilterLiteral:
		return e.evalFilterLiteral(node)

	default:
		return object.NULL
	}
}

func (e *Evaluator) evalTagStatement(node *ast.TagStatement) object.Object {
	var results []object.Object

	for _, node := range node.Nodes {
		switch node := node.(type) {
		case *ast.Identifier:
			results = append(results, object.New(node.Value))
		default:
			results = append(results, e.Eval(node))
		}
	}

	var blockStmts []tag.Statement
	if node.Tag.Block() {
		for _, stmt := range node.BlockStatement.Statements {
			blockStmts = append(blockStmts, stmt.(tag.Statement))
		}
	}

	return node.Tag.Eval(e, results, blockStmts)
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJ_NUMBER && right.Type() == object.OBJ_NUMBER:
		return e.evalNumberOperation(operator, left, right)
	case operator == "==":
		return convertBoolean(left == right)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberOperation(operator string, left, right object.Object) object.Object {
	leftVal := left.Value().(float64)
	rightVal := right.Value().(float64)

	switch operator {
	case "+":
		return object.New(leftVal + rightVal)
	case "-":
		return object.New(leftVal - rightVal)
	case "*":
		return object.New(leftVal * rightVal)
	case "/":
		return object.New(leftVal / rightVal)
	case ">":
		return convertBoolean(leftVal > rightVal)
	case "<":
		return convertBoolean(leftVal < rightVal)
	case ">=":
		return convertBoolean(leftVal >= rightVal)
	case "<=":
		return convertBoolean(leftVal <= rightVal)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalPrefix(operator string, right object.Object) object.Object {
	switch {
	case right.Type() == object.OBJ_NUMBER:
		return e.evalNumberPrefix(operator, right)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return object.New(right.Value().(float64) * -1)
	default:
		return right
	}
}

func (e *Evaluator) evalIdentifier(name string) object.Object {
	return e.context.Get(name)
}

func (e *Evaluator) evalFilterLiteral(node *ast.FilterLiteral) object.Object {
	filterObj := &object.Filter{
		Name:       node.Name,
		Parameters: make(map[string]object.Object),
	}

	for paramName, paramExp := range node.Parameters {
		filterObj.Parameters[paramName] = e.Eval(paramExp)
	}

	return filterObj
}

func (e *Evaluator) evalFilter(input, filter object.Object) object.Object {
	filterName := filter.(*object.Filter).Name
	filterParams := filter.(*object.Filter).Parameters

	filterFunc := late.FindFilter(filterName)

	if filterFunc == nil {
		return object.NULL
	}

	return filterFunc.Call(input, filterParams)
}

func convertBoolean(value bool) object.Object {
	if value {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
